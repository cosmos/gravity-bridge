use contact::client::Contact;
use std::time::Duration;
use std::time::Instant;
use tokio::time::delay_for;

pub async fn wait_for_next_cosmos_block(contact: &Contact, timeout: Duration) {
    let start = Instant::now();
    let mut last_height = None;
    while Instant::now() - start < timeout {
        if let Ok(block_response) = contact.get_latest_block().await {
            if let Some(block) = block_response.block {
                if last_height.is_some() && block.last_commit.height > last_height.unwrap() {
                    return;
                }
                last_height = Some(block.last_commit.height);
            }
        }
        delay_for(Duration::from_secs(1)).await;
    }
    panic!("Cosmos chain took too long to produce a block!");
}

pub async fn wait_for_cosmos_online(contact: &Contact, timeout: Duration) {
    // we have a block now, wait for a few more.
    wait_for_next_cosmos_block(contact, timeout).await;
    wait_for_next_cosmos_block(contact, timeout).await;
    wait_for_next_cosmos_block(contact, timeout).await;
}
