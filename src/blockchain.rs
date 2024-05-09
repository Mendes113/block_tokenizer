


use crate::block::Block;

pub struct Blockchain {
    blocks: Vec<Block>,
}

impl Blockchain {

    //nova blockchain
    pub fn new() -> Blockchain {
        let genesis_block = Block::new(0, String::from("Genesis Block"), String::from("0"));
        Blockchain {
            blocks: vec![genesis_block],
        }
    }

    //adicionar um bloco à blockchain
    pub fn add_block(&mut self, data: String) {
        let previous_hash = self.blocks.last().unwrap().hash.clone();
        let new_block = Block::new(self.blocks.len() as u32, data, previous_hash);
        self.blocks.push(new_block);
    }


    //obter os blocos da blockchain
    pub fn get_blocks(&self) -> &Vec<Block> {
        &self.blocks
    }

    //obter o índice do próximo bloco
    pub fn get_next_index(&self) -> u32 {
        self.blocks.len() as u32
    }

    //obter o hash do último bloco
    pub fn get_last_hash(&self) -> String {
        self.blocks.last().unwrap().hash.clone()
    }
    



}
