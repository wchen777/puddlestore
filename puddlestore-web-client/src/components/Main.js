import React, { useState } from 'react';
import {
  Box,
  Text,
  Link,
  VStack,
  Code,
  Grid,
  HStack,
  Heading,
  Select,
  Badge,
  Button,
  AlertDialog,
  AlertDialogBody,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogContent,
  AlertDialogOverlay,
  useDisclosure,
  Divider,
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
import ControlPanel from './ControlPanel.js';
import CommandHistory from './CommandHistory.js';
import OpenData from './OpenData.js';

const client = new PuddleStoreClient('http://localhost:3334'); // connects to the envoy proxy

export default function Main() {
  const [connected, setConnected] = useState(false); // is the client connected to the server?
  const [clientData, setClientData] = useState({}); // client data, such as client ID
  const [openFDs, setFDs] = useState([]); // open file descriptors
  const [commandHistory, setCommandHistory] = useState([]); // command history

  const [error, setError] = useState(null);

  // ------------------------ HANDLERS FOR CLIENT REQUESTS ------------------------ //

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
        // wipe all client data
        setConnected(false);
        setClientData({});
        setFDs([]);
        setCommandHistory([]);
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
          setFDs([
            ...openFDs,
            { fd: openResp.getFd(), filepath: filepath, write },
          ]);
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
          setCommandHistory([
            ...commandHistory,
            {
              command: 'close',
              argString: `fd: ${fd}`,
              result: 'success',
            },
          ]);
          const newFDs = openFDs.filter(fdObj => fdObj.fd !== parseInt(fd));
          console.log('newFDs', newFDs);
          setFDs(newFDs);
        }
      }
    );
  };

  const readRequestPuddleStore = ({ fd, offset, size }) => {
    let utf8Decode = new TextDecoder();
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
              result: utf8Decode.decode(readResp.getData()),
            },
          ]);
        }
      }
    );
  };

  const writeRequestPuddleStore = ({ fd, data, offset }) => {
    let utf8Encode = new TextEncoder(); // convert data string to byte array for write

    client.clientWrite(
      new WriteMessage([
        clientData.clientID,
        fd,
        utf8Encode.encode(data),
        offset,
      ]),
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
              result: listResp.array[0].toString(),
            },
          ]);
        }
      }
    );
  };

  console.log(openFDs);

  // ------------------------------------------------------------------------ //

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
        <VStack spacing={6} py={20}>
          <VStack mt={8}>
            <Heading as="h2" size="xl">
              {' '}
              Welcome to PuddleStore!{' '}
            </Heading>
            <Heading as="h4" size="md">
              {' '}
              Your ID is {clientData.clientID}{' '}
            </Heading>
          </VStack>

          <Divider orientation="horizontal" />
          <HStack spacing={5}>
            <CommandHistory commandHistory={commandHistory} />
            <OpenData openFDs={openFDs} />
          </HStack>

          <Divider orientation="horizontal" />
          <ControlPanel
            open={openRequestPuddleStore}
            close={closeRequestPuddleStore}
            read={readRequestPuddleStore}
            write={writeRequestPuddleStore}
            mkdir={mkdirRequestPuddleStore}
            remove={removeRequestPuddleStore}
            list={listRequestPuddleStore}
          />
          <Divider orientation="horizontal" />
          <HStack spacing={5}>
            <Button
              colorScheme="purple"
              onClick={() => disconnectFromPuddleStore()}
              size="lg"
            >
              Exit PuddleStore
            </Button>
            <Badge colorScheme={'red'}>
              {'<-'} REMEMBER TO CLICK THIS BUTTON BEFORE EXITING! <br />
              Puddlestore server does not remove clients automatically yet.
            </Badge>
          </HStack>
        </VStack>
      )}

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
