//! This is the testing module for arbitrary logic functionality. This is where instead of managing transfers directly the bridge simply passes an
//! arbitrary call to an arbitrary sub contract along with a specific amount of funds, allowing for execution of whatever command is required

use crate::TOTAL_TIMEOUT;
use deep_space::Contact;
use gravity_proto::gravity::query_client::QueryClient as GravityQueryClient;
use tokio::time::sleep as delay_for;
use tonic::transport::Channel;
use web30::client::Web3;

pub async fn arbitrary_logic_test(
    _web30: &Web3,
    _grpc_client: GravityQueryClient<Channel>,
    _contact: &Contact,
) {
    delay_for(TOTAL_TIMEOUT).await;
}
