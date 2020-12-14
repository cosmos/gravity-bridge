use contact::client::Contact;
use std::time::Duration;
use std::time::Instant;
use tokio::time::delay_for;

pub async fn wait_for_next_cosmos_block(contact: &Contact, timeout: Duration) {
    let start = Instant::now();
    let current_block = contact
        .get_latest_block()
        .await
        .unwrap()
        .block
        .last_commit
        .height;
    while current_block
        == contact
            .get_latest_block()
            .await
            .unwrap()
            .block
            .last_commit
            .height
    {
        delay_for(Duration::from_secs(1)).await;
        if Instant::now() - start > timeout {
            panic!("Cosmos chain took too long to produce a block!");
        }
    }
}

pub async fn wait_for_cosmos_online(contact: &Contact, timeout: Duration) {
    let start = Instant::now();
    let mut current_block = contact.get_latest_block().await;
    while current_block.is_err() {
        delay_for(Duration::from_secs(1)).await;
        current_block = contact.get_latest_block().await;
        if Instant::now() - start > timeout {
            panic!("Cosmos chain took too long to start!");
        }
    }
    // we have a block now, wait for a few more.
    wait_for_next_cosmos_block(contact, timeout).await;
    wait_for_next_cosmos_block(contact, timeout).await;
}
