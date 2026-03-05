use serde::{Deserialize, Serialize};
use std::collections::HashSet;

/// Manages cluster membership and node registration
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ClusterMembership {
    /// Current broker ID
    broker_id: u32,
    
    /// Set of active broker IDs
    active_brokers: HashSet<u32>,
}

impl ClusterMembership {
    pub fn new(broker_id: u32) -> Self {
        let mut active_brokers = HashSet::new();
        active_brokers.insert(broker_id);
        
        Self {
            broker_id,
            active_brokers,
        }
    }
    
    pub fn broker_id(&self) -> u32 {
        self.broker_id
    }
    
    pub fn add_broker(&mut self, broker_id: u32) {
        self.active_brokers.insert(broker_id);
    }
    
    pub fn remove_broker(&mut self, broker_id: u32) {
        self.active_brokers.remove(&broker_id);
    }
    
    pub fn active_brokers(&self) -> Vec<u32> {
        self.active_brokers.iter().copied().collect()
    }
    
    pub fn is_leader(&self, leader_id: u32) -> bool {
        leader_id == self.broker_id
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    
    #[test]
    fn test_cluster_membership_creation() {
        let membership = ClusterMembership::new(1);
        assert_eq!(membership.broker_id(), 1);
    }
    
    #[test]
    fn test_add_remove_brokers() {
        let mut membership = ClusterMembership::new(1);
        membership.add_broker(2);
        membership.add_broker(3);
        assert_eq!(membership.active_brokers().len(), 3);
        
        membership.remove_broker(2);
        assert_eq!(membership.active_brokers().len(), 2);
    }
}
