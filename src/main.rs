mod blockchain;
mod block;

use block::Block;
use blockchain::Blockchain;

fn main() {
    // Criando uma nova blockchain
    let mut blockchain = Blockchain::new();

    // Adicionando alguns blocos à blockchain com diferentes níveis de dificuldade
    let difficulties = vec![2, 3, 4];
    for difficulty in difficulties {
        let data = format!("Dados do Bloco com dificuldade {}", difficulty);
        let _new_block = Block::create_block(blockchain.get_next_index(), data, blockchain.get_last_hash(), difficulty);
        blockchain.add_block("Block".to_owned());
    }

    // Exibindo os blocos da blockchain
    for block in blockchain.get_blocks() {
        println!("Index: {}", block.index);
        println!("Timestamp: {}", block.timestamp);
        println!("Data: {}", block.data);
        println!("Previous Hash: {}", block.previous_hash);
        println!("Hash: {}", block.hash);
        println!();
    }


    // Mineração de um bloco com dificuldade 5
    let data = "Dados do Bloco com dificuldade 5".to_owned();
    let mut new_block = Block::create_block(blockchain.get_next_index(), data, blockchain.get_last_hash(), 5);
    new_block.mine_block_with_difficulty(5);
    blockchain.add_block("Block".to_owned());
    //mine and print
    new_block.mine_block_with_difficulty_and_print(5);
}
