import React, { useState } from 'react';
import {
  Box,
  Text,
  Link,
  VStack,
  Code,
  Grid,
  HStack,
  Select,
  Button,
  AlertDialog,
  AlertDialogBody,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogContent,
  AlertDialogOverlay,
  useDisclosure,
} from '@chakra-ui/react';

import {
  OpenMessage,
  CloseMessage,
  ReadMessage,
  WriteMessage,
  MkdirMessage,
  RemoveMessage,
  ListMessage,
  ReadResponse,
  ListResponse,
  OpenResponse,
  Success,
  ClientID,
  Empty,
} from '../puddlestore/puddlestore_pb.js';

import { PuddleStoreClient } from '../puddlestore/puddlestore_grpc_web_pb.js';
import './Main.css';

const client = new PuddleStoreClient('http://localhost:3334'); // connects to the envoy proxy

export default function Main() {
  const [connected, setConnected] = useState(false);
  const [clientData, setClientData] = useState({});
  const [error, setError] = useState(null);

  const connectToPuddleStore = () => {
    client.clientConnect(new Empty(), {}, (err, id) => {
      if (err) {
        console.log(err);
        setError(err.toString());
        onOpen();
      } else {
        setConnected(true);
        setClientData({ ...clientData, clientID: id.array[0] });
        console.log('client id resp', id.array[0]);
      }
    });
  };

  const disconnectFromPuddleStore = () => {
    client.clientExit(new ClientID([clientData.clientID]), {}, (err, _) => {
      if (err) {
        console.log(err);
        setError(err.toString());
        onOpen(); // open error modal
      } else {
        setConnected(false);
        setClientData({});
      }
    });
  };

  const { isOpen, onOpen, onClose } = useDisclosure();
  const cancelRef = React.useRef();

  return (
    <Box fontSize="l" className="box">
      {!connected && (
        <Button
          colorScheme="purple"
          onClick={() => connectToPuddleStore()}
          size="lg"
        >
          Connect to PuddleStore
        </Button>
      )}

      {connected && (
        <Button
          colorScheme="purple"
          onClick={() => disconnectFromPuddleStore()}
          size="lg"
        >
          Exit PuddleStore{'\n'}
          (Remember to click this! Puddlestore server does not remove clients
          automatically yet)
        </Button>
      )}

      {/* <HStack spacing={4}>

            <Select placeholder='Command'>
                <option value='Open'>Open</option>
                <option value='Close'>Close</option>
                <option value='Read'>Read</option>
                <option value='Write'>Write</option>
                <option value='Mkdir'>Mkdir</option>
                <option value='Remove'>Remove</option>
                <option value='List'>List</option>
            </Select>
        </HStack> */}

      <AlertDialog
        isOpen={isOpen}
        leastDestructiveRef={cancelRef}
        onClose={onClose}
      >
        <AlertDialogOverlay>
          <AlertDialogContent>
            <AlertDialogHeader fontSize="lg" fontWeight="bold">
              PuddleStore Error
            </AlertDialogHeader>

            <AlertDialogBody>{error}</AlertDialogBody>

            <AlertDialogFooter>
              <Button ref={cancelRef} onClick={onClose}>
                Close
              </Button>
            </AlertDialogFooter>
          </AlertDialogContent>
        </AlertDialogOverlay>
      </AlertDialog>
    </Box>
  );
}
