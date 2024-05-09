extern crate chrono;
use chrono::prelude::*;
use crypto::digest::Digest;
use crypto::sha2::Sha256;



//uma blockchain é uma lista de blocos
//cada bloco contém um índice, um timestamp, dados, o hash do bloco anterior e o hash do bloco atual
pub struct Block {
    pub index: u32,
    pub timestamp: i64,
    pub data: String,
    pub previous_hash: String,
    pub hash: String,
    pub nonce : u32,
}

impl Block {

    //novo bloco
    pub fn new(index: u32, data: String, previous_hash: String) -> Block {
        let timestamp = Utc::now().timestamp();
        let hash = Block::calculate_hash(index, timestamp, &data, &previous_hash, 0);

        Block {
            index,
            timestamp,
            data,
            previous_hash,
            hash,
            nonce: 0,
        }
    }

    fn calculate_hash(index: u32, timestamp: i64, data: &str, previous_hash: &str, nonce: u64) -> String {
        let mut hasher = Sha256::new();
        hasher.input_str(&index.to_string());
        hasher.input_str(&timestamp.to_string());
        hasher.input_str(data);
        hasher.input_str(previous_hash);
        hasher.input_str(&nonce.to_string());
        hasher.result_str()
    }



    fn is_valid_hash(hash: &str, difficulty: usize) -> bool {
        hash.starts_with(&"0".repeat(difficulty))
    }

    fn mine_block(block: &mut Block, difficulty: usize) {
        while !Block::is_valid_hash(&block.hash, difficulty) {
            block.nonce += 1;
            block.hash = Block::calculate_hash(block.index, block.timestamp, &block.data, &block.previous_hash, block.nonce as u64);
        }
    }


    pub fn mine_block_with_difficulty(&mut self, difficulty: usize) {
        Block::mine_block(self, difficulty);
    }

    pub fn mine_block_with_difficulty_and_print(&mut self, difficulty: usize) {
        Block::mine_block(self, difficulty);
        println!("Block mined: {}", self.hash);
    }

    pub fn create_block(index: u32, data: String, previous_hash: String, difficulty: usize) -> Block {
        let timestamp = Utc::now().timestamp();
        let mut block = Block {
            index,
            timestamp,
            data,
            previous_hash,
            hash: String::new(),
            nonce: 0,
        };
        block.mine_block_with_difficulty(difficulty);
        block
    }





}
