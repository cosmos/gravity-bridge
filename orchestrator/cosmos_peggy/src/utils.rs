use contact::client::Contact;
use std::time::Duration;
use std::{thread, time::Instant};

pub async fn wait_for_next_cosmos_block(contact: &Contact) {
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
        thread::sleep(Duration::from_secs(1))
    }
}

pub async fn wait_for_cosmos_online(contact: &Contact, timeout: Duration) {
    let start = Instant::now();
    let mut current_block = contact.get_latest_block().await;
    while current_block.is_err() {
        thread::sleep(Duration::from_secs(1));
        current_block = contact.get_latest_block().await;
        if Instant::now() - start > timeout {
            panic!("Cosmos chain took too long to start!");
        }
    }
}
