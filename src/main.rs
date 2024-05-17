mod blockchain;
mod block;

use block::Block;
use blockchain::Blockchain;

fn main() {
    
    let mut blockchain = Blockchain::new();

    
    let difficulties = vec![2, 3, 4];
    for difficulty in difficulties {
        let data = format!("Dados do Bloco com dificuldade {}", difficulty);
        blockchain.add_block(data.clone(), difficulty); 
    }

    
    for block in blockchain.get_blocks() {
        println!("Index: {}", block.index);
        println!("Timestamp: {}", block.timestamp);
        println!("Data: {}", block.data);
        println!("Previous Hash: {}", block.previous_hash);
        println!("Hash: {}", block.hash);
        println!();
    }

    
    let data_fail = "Dados do Bloco com dificuldade 5 que falhará na mineração".to_owned();
    let mut new_block_fail = Block::create_block(blockchain.get_next_index(), data_fail.clone(), blockchain.get_last_hash(), 5);
    new_block_fail.mine_block_with_difficulty();
    println!("Resultado da mineração do bloco com dificuldade 5 que falha:");
    new_block_fail.mine_block_with_difficulty_and_print();
    println!();

    
    let data_success = "Dados do Bloco com dificuldade 4 que será minerado com sucesso".to_owned();
    let mut new_block_success = Block::create_block(blockchain.get_next_index(), data_success.clone(), blockchain.get_last_hash(), 4);
    new_block_success.mine_block_with_difficulty();
    println!("Resultado da mineração do bloco com dificuldade 4 que terá sucesso:");
    new_block_success.mine_block_with_difficulty_and_print();
    println!();

    
    blockchain.add_block(data_success, 4);

    
    for block in blockchain.get_blocks() {
        println!("Index: {}", block.index);
        println!("Timestamp: {}", block.timestamp);
        println!("Data: {}", block.data);
        println!("Previous Hash: {}", block.previous_hash);
        println!("Hash: {}", block.hash);
        println!();
    }

     
     let data_fail = "Dados do Bloco com dificuldade 6 que falhará na mineração".to_owned();
     let mut new_block_fail = Block::create_block(blockchain.get_next_index(), data_fail.clone(), blockchain.get_last_hash(), 6);
     new_block_fail.mine_block_with_difficulty();
     println!("Resultado da mineração do bloco com dificuldade 6 que falha:");
     new_block_fail.mine_block_with_difficulty_and_print();
     println!();
     
}
