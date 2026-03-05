pub mod manager;
pub mod offset;

pub use manager::PartitionManager;
pub use offset::OffsetManager;

/// Partition state and management
#[derive(Debug)]
pub struct Partition {
    id: u32,
    topic: String,
    manager: PartitionManager,
    offset_manager: OffsetManager,
}

impl Partition {
    pub fn new(id: u32, topic: String) -> Self {
        Self {
            id,
            topic,
            manager: PartitionManager::new(id),
            offset_manager: OffsetManager::new(),
        }
    }
}
