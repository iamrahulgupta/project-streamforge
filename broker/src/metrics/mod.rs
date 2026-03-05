//! Metrics and monitoring for the broker
//! 
//! Provides Prometheus metrics, health checks, and observability hooks

use metrics::{counter, gauge, histogram};

/// Metrics for broker operations
#[derive(Debug)]
pub struct BrokerMetrics {
    /// Number of messages produced
    pub messages_produced: u64,
    
    /// Number of messages consumed
    pub messages_consumed: u64,
    
    /// Number of bytes written
    pub bytes_written: u64,
    
    /// Number of bytes read
    pub bytes_read: u64,
}

impl BrokerMetrics {
    pub fn new() -> Self {
        Self {
            messages_produced: 0,
            messages_consumed: 0,
            bytes_written: 0,
            bytes_read: 0,
        }
    }
    
    /// Record a message production
    pub fn record_produce(&mut self, bytes: u64) {
        self.messages_produced = self.messages_produced.saturating_add(1);
        self.bytes_written = self.bytes_written.saturating_add(bytes);
        
        counter!("broker_messages_produced", 1);
        counter!("broker_bytes_written", bytes);
    }
    
    /// Record a message consumption
    pub fn record_consume(&mut self, bytes: u64) {
        self.messages_consumed = self.messages_consumed.saturating_add(1);
        self.bytes_read = self.bytes_read.saturating_add(bytes);
        
        counter!("broker_messages_consumed", 1);
        counter!("broker_bytes_read", bytes);
    }
    
    /// Record operation latency
    pub fn record_latency(&self, operation: &str, duration_ms: u64) {
        histogram!("broker_operation_latency_ms", duration_ms as f64);
    }
    
    /// Update active connections gauge
    pub fn set_active_connections(&self, count: u64) {
        gauge!("broker_active_connections", count as f64);
    }
}

impl Default for BrokerMetrics {
    fn default() -> Self {
        Self::new()
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    
    #[test]
    fn test_metrics_creation() {
        let metrics = BrokerMetrics::new();
        assert_eq!(metrics.messages_produced, 0);
        assert_eq!(metrics.messages_consumed, 0);
    }
    
    #[test]
    fn test_record_produce() {
        let mut metrics = BrokerMetrics::new();
        metrics.record_produce(1024);
        
        assert_eq!(metrics.messages_produced, 1);
        assert_eq!(metrics.bytes_written, 1024);
    }
    
    #[test]
    fn test_record_consume() {
        let mut metrics = BrokerMetrics::new();
        metrics.record_consume(512);
        
        assert_eq!(metrics.messages_consumed, 1);
        assert_eq!(metrics.bytes_read, 512);
    }
}
