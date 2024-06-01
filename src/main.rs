#[macro_use]
extern crate lazy_static;
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
use log::{error, info, LevelFilter};

use block::Block;
use blockchain::Blockchain;

lazy_static! {
    static ref BLOCKCHAIN: Mutex<Blockchain> = Mutex::new(Blockchain::new("blockchain_data.json"));
}

fn main() {
    env_logger::Builder::new()
        .filter_level(LevelFilter::Info)
        .init();
    info!("Starting Rust blockchain...");

    let consumer: BaseConsumer = ClientConfig::new()
        .set("group.id", "rust_blockchain_group")
        .set("bootstrap.servers", "127.0.0.1:9092")
        .create()
        .expect("Consumer creation failed");

    consumer.subscribe(&["blockchain"]).expect("Failed to subscribe to topic");
    info!("Subscribed to blockchain topic");

    loop {
        match consumer.poll(Duration::from_millis(10)) {
            Some(Ok(message)) => {
                if let Some(payload) = message.payload_view::<str>() {
                    match payload {
                        Ok(text) => {
                            handle_message(text);
                        },
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
    info!("Received message: {}", message); // Adiciona log de entrada

    let parts: Vec<&str> = message.split(':').collect();
    if parts.len() == 2 && parts[0] == "add_block" {
        let data = parts[1].to_string();
        let difficulty = 2;  // Define a dificuldade conforme necessário

        let mut blockchain = BLOCKCHAIN.lock().unwrap();
        match blockchain.add_block(data, difficulty) {
            Ok(_) => {
                info!("Block added successfully");
                let last_block_json = blockchain.get_last_block_to_json().unwrap();
                // Configurar o Kafka Producer
                let producer: BaseProducer = ClientConfig::new()
                    .set("bootstrap.servers", "127.0.0.1:9092")
                    .create()
                    .expect("Producer creation failed");

                // Enviar mensagem para o Kafka
                info!("Sending message to blockchain_data topic");
                let result = producer.send(
                    BaseRecord::to("blockchain_data")
                        .key("blockchain_data")
                        .payload(&last_block_json)
                );

                match result {
                    Ok(_) => {
                        info!("Message sent successfully");
                        producer.poll(Duration::from_millis(10));
                    },
                    Err(e) => error!("Failed to send message: {:?}", e),
                }
            },
            Err(e) => {
                error!("Failed to add block: {:?}", e);
                // Configurar o Kafka Producer para mensagens de erro
                let producer: BaseProducer = ClientConfig::new()
                    .set("bootstrap.servers", "127.0.0.1:9092")
                    .create()
                    .expect("Producer creation failed");

                // Enviar mensagem de erro para o Kafka
                info!("Sending error message to blockchain_error topic");
                let error_message = format!("Failed to add block: {:?}", e);
                let result = producer.send(
                    BaseRecord::to("blockchain_err")
                        .key("blockchain_error")
                        .payload(&error_message)
                );

                match result {
                    Ok(_) => {
                        info!("Error message sent successfully");
                        producer.poll(Duration::from_millis(10));
                    },
                    Err(e) => error!("Failed to send error message: {:?}", e),
                }
            },
        }
    }
}
