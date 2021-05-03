use deep_space::Contact;
use std::time::Duration;

pub async fn wait_for_cosmos_online(contact: &Contact, timeout: Duration) {
    // we have a block now, wait for a few more.
    contact.wait_for_next_block(timeout).await.unwrap();
    contact.wait_for_next_block(timeout).await.unwrap();
    contact.wait_for_next_block(timeout).await.unwrap();
}
