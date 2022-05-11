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
  const [connected, setConnected] = useState(false); // is the client connected to the server?
  const [clientData, setClientData] = useState({}); // client data, such as client ID
  const [openFDs, setFDs] = useState([]); // open file descriptors
  const [commandHistory, setCommandHistory] = useState([]); // command history

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

  const openRequestPuddleStore = ({ filepath, create, write }) => {
    client.clientOpen(
      new OpenMessage([clientData.clientID, filepath, create, write]),
      {},
      (err, openResp) => {
        if (err) {
          console.log(err);
          setError(err.toString());
          onOpen(); // open error modal
        } else {
          console.log('open response', openResp);
          setFDs([...openFDs, { fd: openResp.getFd(), filepath: filepath }]);
        }
      }
    );
  };

  const closeRequestPuddleStore = ({ fd }) => {
    client.clientClose(
      new CloseMessage([clientData.clientID, fd]),
      {},
      (err, _) => {
        if (err) {
          console.log(err);
          setError(err.toString());
          onOpen(); // open error modal
        } else {
          setFDs(openFDs.filter(file => file.fd !== fd));
        }
      }
    );
  };

  const readRequestPuddleStore = ({ fd, offset, size }) => {
    client.clientRead(
      new ReadMessage([clientData.clientID, fd, offset, size]),
      {},
      (err, readResp) => {
        if (err) {
          console.log(err);
          setError(err.toString());
          onOpen(); // open error modal
        } else {
          console.log('read response', readResp);

          // handle read response here!!
          /* steps:
            - create new command result component
            - containing: data from command, and result of command
          */
          setCommandHistory([
            ...commandHistory,
            {
              command: 'read',
              argString: `fd: ${fd}, offset: ${offset}, size: ${size}`,
              result: readResp.getData().toString(),
            },
          ]);
        }
      }
    );
  };

  const writeRequestPuddleStore = ({ fd, data, offset }) => {
    client.clientWrite(
      new WriteMessage([clientData.clientID, fd, data, offset]),
      {},
      (err, _) => {
        if (err) {
          console.log(err);
          setError(err.toString());
          onOpen(); // open error modal
        } else {
          // create a component on success
          setCommandHistory([
            ...commandHistory,
            {
              command: 'write',
              argString: `fd: ${fd}, data: ${data}, offset: ${offset}`,
              result: 'success',
            },
          ]);
        }
      }
    );
  };

  const mkdirRequestPuddleStore = ({ path }) => {
    client.clientMkdir(
      new MkdirMessage([clientData.clientID, path]),
      {},
      (err, _) => {
        if (err) {
          console.log(err);
          setError(err.toString());
          onOpen(); // open error modal
        } else {
          // create a component on success
          setCommandHistory([
            ...commandHistory,
            {
              command: 'mkdir',
              argString: `path: ${path}`,
              result: 'success',
            },
          ]);
        }
      }
    );
  };

  const removeRequestPuddleStore = ({ path }) => {
    client.clientRemove(
      new RemoveMessage([clientData.clientID, path]),
      {},
      (err, _) => {
        if (err) {
          console.log(err);
          setError(err.toString());
          onOpen(); // open error modal
        } else {
          // create a component on success
          setCommandHistory([
            ...commandHistory,
            {
              command: 'remove',
              argString: `path: ${path}`,
              result: 'success',
            },
          ]);
        }
      }
    );
  };

  const listRequestPuddleStore = ({ path }) => {
    client.clientList(
      new ListMessage([clientData.clientID, path]),
      {},
      (err, listResp) => {
        if (err) {
          console.log(err);
          setError(err.toString());
          onOpen(); // open error modal
        } else {
          console.log('list response', listResp);
          // handle list response here!!
          /* steps:
                    - create new command result component
                    - containing: data from command, and result of command
                */
          setCommandHistory([
            ...commandHistory,
            {
              command: 'list',
              argString: `path: ${path}`,
              result: listResp.getData().toString(),
            },
          ]);
        }
      }
    );
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
        <HStack spacing={5}>
          <Button
            colorScheme="purple"
            onClick={() => disconnectFromPuddleStore()}
            size="lg"
          >
            Exit PuddleStore
          </Button>
          <Text fontSize="xl">
            (REMEMBER TO CLICK THIS BUTTON! <br />
            Puddlestore server does not remove clients automatically yet.)
          </Text>
        </HStack>
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
