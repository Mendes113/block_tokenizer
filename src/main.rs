#[macro_use] extern crate lazy_static;
extern crate chrono;
extern crate crypto;
extern crate serde;
extern crate serde_json;
extern crate rdkafka;

mod block;
mod blockchain;

use std::sync::Mutex;
use std::time::Duration;
use rdkafka::config::ClientConfig;
use rdkafka::consumer::{BaseConsumer, Consumer};
use rdkafka::message::Message;
use rdkafka::producer::{BaseProducer, BaseRecord};
use log::{info, error};

use block::Block;
use blockchain::Blockchain;

lazy_static! {
    static ref BLOCKCHAIN: Mutex<Blockchain> = Mutex::new(Blockchain::new("blockchain_data.json"));
}

fn main() {
    env_logger::init();
    info!("Starting Rust blockchain...");

    let consumer: BaseConsumer = ClientConfig::new()
        .set("group.id", "rust_blockchain_group")
        .set("bootstrap.servers", "localhost:9092")
        .create()
        .expect("Consumer creation failed");

    consumer.subscribe(&["blockchain"]).expect("Failed to subscribe to topic");

    loop {
        match consumer.poll(Duration::from_millis(100)) {
            Some(Ok(message)) => {
                if let Some(payload) = message.payload_view::<str>() {
                    match payload {
                        Ok(text) => handle_message(text),
                        Err(e) => error!("Error while deserializing message payload: {:?}", e),
                    }
                }
            },
            Some(Err(e)) => error!("Error while receiving message: {:?}", e),
            None => (),
        }
    }
}

fn handle_message(message: &str) {
    let parts: Vec<&str> = message.split(':').collect();
    if parts.len() == 2 && parts[0] == "add_block" {
        let data = parts[1].to_string();
        let difficulty = 2;  // Define a dificuldade conforme necessário

        let mut blockchain = BLOCKCHAIN.lock().unwrap();
        match blockchain.add_block(data, difficulty) {
            Ok(_) => info!("Block added successfully"),
            Err(e) => error!("Failed to add block: {:?}", e),
        }
    }
}
