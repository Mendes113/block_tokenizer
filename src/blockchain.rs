use crate::block::Block;
use log::{error, info};
use serde_json;
use std::fs;
use std::path::Path;

pub struct Blockchain {
    blocks: Vec<Block>,
    data_file: String,
}

impl Blockchain {
    pub fn new(data_file: &str) -> Blockchain {
        let mut blockchain = Blockchain {
            blocks: Vec::new(),
            data_file: data_file.to_string(),
        };
        blockchain.load_blocks();
        if blockchain.blocks.is_empty() {
            let genesis_block = Block::new(0, String::from("Genesis Block"), String::from("0"), 1);
            blockchain.blocks.push(genesis_block);
        }
        blockchain
    }

    pub fn add_block(&mut self, data: String, difficulty: usize) -> Result<(), &'static str> {
        info!("Adding block to blockchain");

        if !self.is_data_valid(&data) {
            error!("Invalid data");
            return Err("Invalid data");
        }

        let previous_hash = self.get_last_hash();
        let new_block = Block::create_block(self.get_next_index(), data, previous_hash, difficulty);

        if !self.is_block_valid(&new_block, difficulty) {
            error!("Block is not valid");
            return Err("Block is not valid");
        }

        self.blocks.push(new_block);
        self.save_blocks();
        Ok(())
    }

    fn is_block_valid(&self, new_block: &Block, difficulty: usize) -> bool {
        if let Some(last_block) = self.blocks.last() {
            if !new_block.is_valid(Some(last_block), difficulty) {
                return false;
            }
            if !self.is_timestamp_valid(new_block) {
                return false;
            }
        }
        true
    }

    fn is_timestamp_valid(&self, new_block: &Block) -> bool {
        if let Some(last_block) = self.blocks.last() {
            new_block.timestamp > last_block.timestamp
        } else {
            true
        }
    }

    fn is_data_valid(&self, data: &String) -> bool {
        !data.is_empty()
    }

    pub fn get_last_block_to_json(&self) -> Option<String> {
        self.blocks.last().map(|block| {
            serde_json::to_string(block).expect("Failed to serialize block to JSON")
        })
    }

    pub fn get_blocks(&self) -> &Vec<Block> {
        &self.blocks
    }

    pub fn get_next_index(&self) -> u32 {
        self.blocks.len() as u32
    }

    pub fn get_last_hash(&self) -> String {
        self.blocks.last().unwrap().hash.clone()
    }

    fn save_blocks(&self) {
        let json = serde_json::to_string(&self.blocks).expect("Unable to serialize blocks");
        fs::write(&self.data_file, json).expect("Unable to write blocks to file");
    }

    fn load_blocks(&mut self) {
        if Path::new(&self.data_file).exists() {
            let json = fs::read_to_string(&self.data_file).expect("Unable to read blocks from file");
            self.blocks = serde_json::from_str(&json).expect("Unable to deserialize blocks");
        }
    }
}
