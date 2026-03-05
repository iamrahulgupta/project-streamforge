/// Integration tests for StreamForge Broker
///
/// These tests verify the correct interaction of multiple broker components

#[cfg(test)]
mod integration_tests {
    use streamforge_broker::cluster::membership::ClusterMembership;
    use streamforge_broker::partition::manager::{Message, PartitionManager};
    use streamforge_broker::replication::log_sync::LogEntry;

    #[test]
    fn test_broker_cluster_initialization() {
        let membership = ClusterMembership::new(1);
        assert_eq!(membership.broker_id(), 1);
    }

    #[test]
    fn test_partition_message_flow() {
        let mut partition = PartitionManager::new(0);

        // Append some messages
        for i in 0..5 {
            partition.append(Message {
                offset: i,
                timestamp: 1234567890,
                key: None,
                value: vec![i as u8],
            });
        }

        // Fetch messages
        let messages = partition.get_messages(0, 10);
        assert_eq!(messages.len(), 5);
    }

    #[test]
    fn test_replicated_log_sync() {
        let mut log = streamforge_broker::replication::log_sync::LogSync::new();

        for i in 0..10 {
            log.append_entry(LogEntry {
                offset: i,
                term: 1,
                data: vec![i as u8; 100],
            });
        }

        let entries = log.get_entries(5);
        assert_eq!(entries.len(), 5);
        assert_eq!(entries[0].offset, 5);
    }
}
