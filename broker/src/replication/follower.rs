use tracing::debug;

/// Follower-side replication management
#[derive(Debug)]
pub struct FollowerReplicator {
    broker_id: u32,
    
    /// Current replication offset
    current_offset: u64,
    
    /// Leader broker ID this follower replicates from
    leader_id: Option<u32>,
}

impl FollowerReplicator {
    pub fn new(broker_id: u32) -> Self {
        Self {
            broker_id,
            current_offset: 0,
            leader_id: None,
        }
    }
    
    pub fn broker_id(&self) -> u32 {
        self.broker_id
    }
    
    pub fn current_offset(&self) -> u64 {
        self.current_offset
    }
    
    pub fn set_leader(&mut self, leader_id: u32) {
        self.leader_id = Some(leader_id);
        debug!("Follower {} following leader {}", self.broker_id, leader_id);
    }
    
    /// Update replication progress
    pub fn update_offset(&mut self, offset: u64) {
        self.current_offset = offset;
        debug!(
            "Follower {} updated offset to {}",
            self.broker_id, offset
        );
    }
    
    /// Fetch messages from leader
    pub async fn fetch_from_leader(
        &self,
        partition: u32,
        start_offset: u64,
    ) {
        debug!(
            "Follower {} fetching from leader {:?} (partition {}, offset {})",
            self.broker_id, self.leader_id, partition, start_offset
        );
        // TODO: Implement fetch logic
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    
    #[test]
    fn test_follower_creation() {
        let follower = FollowerReplicator::new(2);
        assert_eq!(follower.broker_id(), 2);
        assert_eq!(follower.current_offset(), 0);
    }
    
    #[test]
    fn test_set_leader() {
        let mut follower = FollowerReplicator::new(2);
        follower.set_leader(1);
        assert_eq!(follower.leader_id, Some(1));
    }
    
    #[test]
    fn test_update_offset() {
        let mut follower = FollowerReplicator::new(2);
        follower.update_offset(50);
        assert_eq!(follower.current_offset(), 50);
    }
}
