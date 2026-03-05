use serde::{Deserialize, Serialize};
use tracing::debug;

/// Log entry for replication
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct LogEntry {
    pub offset: u64,
    pub term: u64,
    pub data: Vec<u8>,
}

/// Synchronizes logs between leader and followers
#[derive(Debug)]
pub struct LogSync {
    entries: Vec<LogEntry>,
}

impl LogSync {
    pub fn new() -> Self {
        Self {
            entries: Vec::new(),
        }
    }
    
    pub fn append_entry(&mut self, entry: LogEntry) {
        debug!("Appending log entry at offset {}", entry.offset);
        self.entries.push(entry);
    }
    
    pub fn get_entries(&self, start_offset: u64) -> Vec<LogEntry> {
        self.entries
            .iter()
            .skip_while(|e| e.offset < start_offset)
            .cloned()
            .collect()
    }
    
    pub fn last_offset(&self) -> Option<u64> {
        self.entries.last().map(|e| e.offset)
    }
    
    pub fn truncate(&mut self, offset: u64) {
        debug!("Truncating log at offset {}", offset);
        self.entries.retain(|e| e.offset < offset);
    }
}

impl Default for LogSync {
    fn default() -> Self {
        Self::new()
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    
    #[test]
    fn test_log_sync_creation() {
        let log = LogSync::new();
        assert!(log.last_offset().is_none());
    }
    
    #[test]
    fn test_append_and_retrieve_entries() {
        let mut log = LogSync::new();
        
        log.append_entry(LogEntry {
            offset: 0,
            term: 1,
            data: vec![1, 2, 3],
        });
        log.append_entry(LogEntry {
            offset: 1,
            term: 1,
            data: vec![4, 5, 6],
        });
        
        assert_eq!(log.last_offset(), Some(1));
        assert_eq!(log.get_entries(0).len(), 2);
        assert_eq!(log.get_entries(1).len(), 1);
    }
    
    #[test]
    fn test_truncate() {
        let mut log = LogSync::new();
        log.append_entry(LogEntry {
            offset: 0,
            term: 1,
            data: vec![],
        });
        log.append_entry(LogEntry {
            offset: 1,
            term: 1,
            data: vec![],
        });
        log.append_entry(LogEntry {
            offset: 2,
            term: 1,
            data: vec![],
        });
        
        log.truncate(2);
        assert_eq!(log.last_offset(), Some(1));
    }
}
