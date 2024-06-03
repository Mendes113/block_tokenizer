use abi::Abi;
use ethers::prelude::*;
use std::convert::TryFrom;
use std::sync::Arc;





pub async fn send_data_to_ethereum(provider_url: &str, contract_address: &str, data: String) -> Result<(), Box<dyn std::error::Error>> {
    // Criar uma variável separada para o provedor
    let provider = Provider::<Http>::try_from(provider_url)?;
    let provider = Arc::new(provider);

    // Criar uma variável separada para o endereço do contrato
    let contract_address = contract_address.parse::<Address>()?;

    // Definir a ABI do contrato
    let abi: Abi = serde_json::from_str(r#"
        [
            {
                "constant": false,
                "inputs": [
                    {
                        "name": "data",
                        "type": "string"
                    }
                ],
                "name": "storeData",
                "outputs": [],
                "payable": false,
                "stateMutability": "nonpayable",
                "type": "function"
            }
        ]
    "#)?;

    // Criar uma variável separada para o contrato
    let contract = Contract::new(contract_address, abi, provider);

    // Criar uma variável separada para a chamada do método
    let method_call = contract.method::<_, H256>("storeData", data)?;

    // Enviar a transação
    let tx = method_call.send().await?;
    println!("Transaction sent: {:?}", tx);

    Ok(())
}
