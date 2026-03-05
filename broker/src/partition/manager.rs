use serde::{Deserialize, Serialize};
use tracing::debug;

/// Message in a partition
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Message {
    pub offset: u64,
    pub timestamp: u64,
    pub key: Option<Vec<u8>>,
    pub value: Vec<u8>,
}

/// Manages a single partition
#[derive(Debug)]
pub struct PartitionManager {
    partition_id: u32,
    
    /// Messages in this partition
    messages: Vec<Message>,
}

impl PartitionManager {
    pub fn new(partition_id: u32) -> Self {
        Self {
            partition_id,
            messages: Vec::new(),
        }
    }
    
    pub fn partition_id(&self) -> u32 {
        self.partition_id
    }
    
    /// Append a message to the partition
    pub fn append(&mut self, message: Message) {
        debug!(
            "Appending message to partition {} at offset {}",
            self.partition_id, message.offset
        );
        self.messages.push(message);
    }
    
    /// Get messages from offset
    pub fn get_messages(&self, start_offset: u64, max_count: u32) -> Vec<Message> {
        self.messages
            .iter()
            .skip_while(|m| m.offset < start_offset)
            .take(max_count as usize)
            .cloned()
            .collect()
    }
    
    /// Get the last offset
    pub fn last_offset(&self) -> Option<u64> {
        self.messages.last().map(|m| m.offset)
    }
    
    /// Get total message count
    pub fn message_count(&self) -> usize {
        self.messages.len()
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    
    #[test]
    fn test_partition_manager_creation() {
        let manager = PartitionManager::new(0);
        assert_eq!(manager.partition_id(), 0);
        assert_eq!(manager.message_count(), 0);
    }
    
    #[test]
    fn test_append_message() {
        let mut manager = PartitionManager::new(0);
        let msg = Message {
            offset: 0,
            timestamp: 1234567890,
            key: Some(vec![1, 2, 3]),
            value: vec![4, 5, 6],
        };
        
        manager.append(msg);
        assert_eq!(manager.message_count(), 1);
        assert_eq!(manager.last_offset(), Some(0));
    }
    
    #[test]
    fn test_get_messages() {
        let mut manager = PartitionManager::new(0);
        for i in 0..5 {
            manager.append(Message {
                offset: i,
                timestamp: 1234567890,
                key: None,
                value: vec![i as u8],
            });
        }
        
        let messages = manager.get_messages(2, 2);
        assert_eq!(messages.len(), 2);
        assert_eq!(messages[0].offset, 2);
    }
}
