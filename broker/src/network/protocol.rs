use serde::{Deserialize, Serialize};

/// Protocol version for inter-broker communication
pub const PROTOCOL_VERSION: u16 = 1;

/// Message types in the broker protocol
#[derive(Debug, Clone, Copy, Serialize, Deserialize, PartialEq, Eq)]
pub enum MessageType {
    Produce = 0,
    Fetch = 1,
    Metadata = 2,
    CommitOffset = 3,
    FetchOffset = 4,
}

impl MessageType {
    pub fn from_code(code: u16) -> Option<Self> {
        match code {
            0 => Some(MessageType::Produce),
            1 => Some(MessageType::Fetch),
            2 => Some(MessageType::Metadata),
            3 => Some(MessageType::CommitOffset),
            4 => Some(MessageType::FetchOffset),
            _ => None,
        }
    }
}

/// Broker protocol message
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ProtocolMessage {
    pub message_type: MessageType,
    pub version: u16,
    pub timestamp: u64,
    pub body: Vec<u8>,
}

/// Broker protocol handler
#[derive(Debug)]
pub struct Protocol {
    version: u16,
}

impl Protocol {
    pub fn new() -> Self {
        Self {
            version: PROTOCOL_VERSION,
        }
    }
    
    pub fn version(&self) -> u16 {
        self.version
    }
    
    /// Serialize a message
    pub fn serialize(&self, message: &ProtocolMessage) -> anyhow::Result<Vec<u8>> {
        Ok(serde_json::to_vec(message)?)
    }
    
    /// Deserialize a message
    pub fn deserialize(&self, data: &[u8]) -> anyhow::Result<ProtocolMessage> {
        Ok(serde_json::from_slice(data)?)
    }
}

impl Default for Protocol {
    fn default() -> Self {
        Self::new()
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    
    #[test]
    fn test_protocol_creation() {
        let protocol = Protocol::new();
        assert_eq!(protocol.version(), PROTOCOL_VERSION);
    }
    
    #[test]
    fn test_message_type_conversion() {
        assert_eq!(MessageType::from_code(0), Some(MessageType::Produce));
        assert_eq!(MessageType::from_code(1), Some(MessageType::Fetch));
        assert_eq!(MessageType::from_code(99), None);
    }
    
    #[test]
    fn test_serialize_deserialize() {
        let protocol = Protocol::new();
        let message = ProtocolMessage {
            message_type: MessageType::Produce,
            version: 1,
            timestamp: 1234567890,
            body: vec![1, 2, 3, 4, 5],
        };
        
        let serialized = protocol.serialize(&message).unwrap();
        let deserialized = protocol.deserialize(&serialized).unwrap();
        
        assert_eq!(message.message_type, deserialized.message_type);
        assert_eq!(message.body, deserialized.body);
    }
}
