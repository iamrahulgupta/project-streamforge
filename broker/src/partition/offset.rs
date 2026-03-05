use std::collections::HashMap;
use tracing::debug;

/// Manages consumer offsets in a partition
#[derive(Debug)]
pub struct OffsetManager {
    /// Consumer group offsets: (consumer_group, consumer_id) -> offset
    offsets: HashMap<(String, String), u64>,
}

impl OffsetManager {
    pub fn new() -> Self {
        Self {
            offsets: HashMap::new(),
        }
    }
    
    /// Commit offset for a consumer
    pub fn commit_offset(
        &mut self,
        consumer_group: String,
        consumer_id: String,
        offset: u64,
    ) {
        debug!(
            "Committing offset {} for consumer {} in group {}",
            offset, consumer_id, consumer_group
        );
        self.offsets.insert((consumer_group, consumer_id), offset);
    }
    
    /// Get offset for a consumer
    pub fn get_offset(&self, consumer_group: &str, consumer_id: &str) -> Option<u64> {
        self.offsets
            .get(&(consumer_group.to_string(), consumer_id.to_string()))
            .copied()
    }
    
    /// Reset offset for a consumer
    pub fn reset_offset(&mut self, consumer_group: &str, consumer_id: &str) {
        self.offsets
            .remove(&(consumer_group.to_string(), consumer_id.to_string()));
    }
}

impl Default for OffsetManager {
    fn default() -> Self {
        Self::new()
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    
    #[test]
    fn test_offset_manager_creation() {
        let manager = OffsetManager::new();
        assert!(manager.get_offset("group1", "consumer1").is_none());
    }
    
    #[test]
    fn test_commit_and_get_offset() {
        let mut manager = OffsetManager::new();
        manager.commit_offset("group1".to_string(), "consumer1".to_string(), 100);
        
        assert_eq!(
            manager.get_offset("group1", "consumer1"),
            Some(100)
        );
    }
    
    #[test]
    fn test_reset_offset() {
        let mut manager = OffsetManager::new();
        manager.commit_offset("group1".to_string(), "consumer1".to_string(), 100);
        manager.reset_offset("group1", "consumer1");
        
        assert!(manager.get_offset("group1", "consumer1").is_none());
    }
}
