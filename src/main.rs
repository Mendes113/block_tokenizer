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
}
