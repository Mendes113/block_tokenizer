use chrono::prelude::*;
use crypto::digest::Digest;
use crypto::sha2::Sha256;
use serde::{Serialize, Deserialize};
use log::info;

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct Block {
    pub index: u32,
    pub timestamp: i64,
    pub data: String,
    pub previous_hash: String,
    pub hash: String,
    pub nonce: u32,
    pub difficulty: usize,
}

impl Block {
    pub fn new(index: u32, data: String, previous_hash: String, difficulty: usize) -> Block {
        let timestamp = Utc::now().timestamp();
        let hash = Block::calculate_hash(index, timestamp, &data, &previous_hash, 0, difficulty);

        Block {
            index,
            timestamp,
            data,
            previous_hash,
            hash,
            nonce: 0,
            difficulty,
        }
    }

    fn calculate_hash(index: u32, timestamp: i64, data: &str, previous_hash: &str, nonce: u64, difficulty: usize) -> String {
        let mut hasher = Sha256::new();
        let input = format!("{}{}{}{}{}", index, timestamp, data, previous_hash, nonce);
        let prefix = "0".repeat(difficulty);

        for nonce_attempt in 0.. {
            let input_with_nonce = format!("{}{}", input, nonce_attempt);
            hasher.input_str(&input_with_nonce);
            let hash = hasher.result_str();
            hasher.reset();

            if hash.starts_with(&prefix) {
                return hash;
            }
        }

        panic!("Unable to find a valid hash after exhausting all nonce values");
    }

    pub fn is_valid_hash(hash: &str, difficulty: usize) -> bool {
        hash.starts_with(&"0".repeat(difficulty))
    }

    fn mine_block(&mut self) {
        while !Block::is_valid_hash(&self.hash, self.difficulty) {
            self.nonce += 1;
            self.hash = Block::calculate_hash(self.index, self.timestamp, &self.data, &self.previous_hash, self.nonce as u64, self.difficulty);
        }
        info!("Block mined: {}", self.hash);
    }

    pub fn create_block(index: u32, data: String, previous_hash: String, difficulty: usize) -> Block {
        let mut block = Block::new(index, data, previous_hash, difficulty);
        block.mine_block();
        block
    }

    pub fn is_valid(&self, previous_block: Option<&Block>, difficulty: usize) -> bool {
        if self.index == 0 {
            return true;
        }
        if let Some(prev) = previous_block {
            if self.index != prev.index + 1 || self.previous_hash != prev.hash {
                return false;
            }
        }
        Block::is_valid_hash(&self.hash, difficulty)
    }
}
