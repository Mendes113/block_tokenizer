


use crate::block::Block;
pub struct Blockchain {
    blocks: Vec<Block>,
}

impl Blockchain {
    pub fn new() -> Blockchain {
        let genesis_block = Block::new(0, String::from("Genesis Block"), String::from("0"), 0);
        Blockchain {
            blocks: vec![genesis_block],
        }
    }

    pub fn add_block(&mut self, data: String, difficulty: usize) -> Result<(), &'static str> {
        
        let previous_hash = match self.blocks.last() {
            Some(last_block) => last_block.hash.clone(),
            None => String::from("0"), 
        };
    
        
        let new_block = Block::create_block(self.blocks.len() as u32, data, previous_hash.clone(), difficulty);
    
        
        if let Some(last_block) = self.blocks.last() {
            if !new_block.is_valid(Some(last_block), difficulty) {
                return Err("Block is not valid");
            }
        } else if !Block::is_valid_hash(&new_block.hash, difficulty) {
            return Err("Block is not valid");
        }
    
        
        self.blocks.push(new_block);
        Ok(())
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
}