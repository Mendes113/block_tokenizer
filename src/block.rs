use chrono::prelude::*;
use crypto::digest::Digest;
use crypto::sha2::Sha256;

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
    let mut hash = String::new();

    
    for nonce_attempt in 0.. {
        let input_with_nonce = format!("{}{}", input, nonce_attempt);
        hasher.input_str(&input_with_nonce);
        hash = hasher.result_str();
        hasher.reset();

        
        if hash.starts_with(&prefix) {
            return hash;
        }
    }
    
    hash
}


    pub fn is_valid_hash(hash: &str, difficulty: usize) -> bool {
        hash.starts_with(&"0".repeat(difficulty))
    }

    fn mine_block(block: &mut Block) {
        while !Block::is_valid_hash(&block.hash, block.difficulty) {
            block.nonce += 1;
            block.hash = Block::calculate_hash(block.index, block.timestamp, &block.data, &block.previous_hash, block.nonce as u64, block.difficulty);
        }
    }

    pub fn mine_block_with_difficulty(&mut self) {
        Block::mine_block(self);
    }

    pub fn mine_block_with_difficulty_and_print(&mut self) {
        Block::mine_block(self);
        println!("Block mined: {}", self.hash);
    }

    pub fn create_block(index: u32, data: String, previous_hash: String, difficulty: usize) -> Block {
        let mut block = Block::new(index, data, previous_hash, difficulty);
        Block::mine_block(&mut block);
        block
    }

    pub fn is_valid(&self, previous_block: Option<&Block>, difficulty: usize) -> bool {
        
        if self.index == 0 {
            
            return true;
        }
        if self.index != previous_block.unwrap().index + 1 {
            return false; 
        }

        
        if self.previous_hash != previous_block.unwrap().hash {
            return false; 
        }

        
        if !Block::is_valid_hash(&self.hash, difficulty) {
            return false; 
        }

        true 
    }
}


