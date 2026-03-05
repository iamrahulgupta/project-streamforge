use serde::{Deserialize, Serialize};
use std::collections::HashMap;

/// Topic metadata information
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct TopicMetadata {
    pub name: String,
    pub partitions: u32,
    pub replication_factor: u16,
}

/// Manages cluster metadata
#[derive(Debug)]
pub struct MetadataManager {
    topics: HashMap<String, TopicMetadata>,
}

impl MetadataManager {
    pub fn new() -> Self {
        Self {
            topics: HashMap::new(),
        }
    }
    
    pub fn add_topic(&mut self, metadata: TopicMetadata) {
        self.topics.insert(metadata.name.clone(), metadata);
    }
    
    pub fn get_topic(&self, name: &str) -> Option<&TopicMetadata> {
        self.topics.get(name)
    }
    
    pub fn remove_topic(&mut self, name: &str) {
        self.topics.remove(name);
    }
    
    pub fn topics(&self) -> Vec<String> {
        self.topics.keys().cloned().collect()
    }
}

impl Default for MetadataManager {
    fn default() -> Self {
        Self::new()
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    
    #[test]
    fn test_metadata_manager_creation() {
        let manager = MetadataManager::new();
        assert_eq!(manager.topics().len(), 0);
    }
    
    #[test]
    fn test_add_topic() {
        let mut manager = MetadataManager::new();
        let metadata = TopicMetadata {
            name: "events".to_string(),
            partitions: 10,
            replication_factor: 3,
        };
        
        manager.add_topic(metadata);
        assert_eq!(manager.topics().len(), 1);
        assert!(manager.get_topic("events").is_some());
    }
}
