import React, { useState } from 'react';
import {
  Box,
  Text,
  VStack,
  HStack,
  Select,
  Button,
  Input,
  Checkbox,
} from '@chakra-ui/react';

export default function ControlPanel({
  open,
  close,
  read,
  write,
  mkdir,
  remove,
  list,
}) {
  const [currentCommand, setCurrentCommand] = useState('Open');

  const [commandData, setCommandData] = useState({});

  const submitCommand = () => {
    switch (currentCommand) {
      case 'Open':
        open(commandData);
        break;
      case 'Close':
        close(commandData);
        break;
      case 'Read':
        read(commandData);
        break;
      case 'Write':
        write(commandData);
        break;
      case 'Mkdir':
        mkdir(commandData);
        break;
      case 'Remove':
        remove(commandData);
        break;
      case 'List':
        list(commandData);
        break;
      default:
        console.log('Invalid command');
        break;
    }
  };

  return (
    <>
      <HStack spacing={4}>
        <Select
          onChange={e => {
            setCurrentCommand(e.target.value);
            setCommandData({});
          }}
          width="auto"
        >
          <option value="Open">Open</option>
          <option value="Close">Close</option>
          <option value="Read">Read</option>
          <option value="Write">Write</option>
          <option value="Mkdir">Mkdir</option>
          <option value="Remove">Remove</option>
          <option value="List">List</option>
        </Select>

        {/* this is scuffed, but it will work */}

        {currentCommand === 'Open' && (
          <>
            <Input
              placeholder="Filepath"
              width="auto"
              onChange={e => {
                setCommandData({ ...commandData, filepath: e.target.value });
              }}
            />

            <Checkbox
              colorScheme="purple"
              onChange={e => {
                setCommandData({ ...commandData, create: e.target.checked });
              }}
            >
              Create
            </Checkbox>

            <Checkbox
              colorScheme="purple"
              onChange={e => {
                setCommandData({ ...commandData, write: e.target.checked });
              }}
            >
              Write
            </Checkbox>
          </>
        )}

        {currentCommand === 'Close' && (
          <>
            <Input
              placeholder="File Descriptor"
              width="auto"
              type={'number'}
              onChange={e => {
                setCommandData({ ...commandData, fd: e.target.value });
              }}
            />
          </>
        )}

        {currentCommand === 'Read' && (
          <>
            <Input
              placeholder="File Descriptor"
              width="auto"
              type={'number'}
              onChange={e => {
                setCommandData({ ...commandData, fd: e.target.value });
              }}
            />

            <Input
              placeholder="Offset"
              width="auto"
              type={'number'}
              onChange={e => {
                setCommandData({ ...commandData, offset: e.target.value });
              }}
            />

            <Input
              placeholder="Size"
              width="auto"
              type={'number'}
              onChange={e => {
                setCommandData({ ...commandData, size: e.target.value });
              }}
            />
          </>
        )}

        {currentCommand === 'Write' && (
          <>
            <Input
              placeholder="File Descriptor"
              width="auto"
              type={'number'}
              onChange={e => {
                setCommandData({ ...commandData, fd: e.target.value });
              }}
            />

            <Input
              placeholder="Data"
              width="auto"
              onChange={e => {
                setCommandData({ ...commandData, data: e.target.value });
              }}
            />

            <Input
              placeholder="Offset"
              width="auto"
              type={'number'}
              onChange={e => {
                setCommandData({ ...commandData, offset: e.target.value });
              }}
            />
          </>
        )}

        {(currentCommand === 'Mkdir' ||
          currentCommand === 'Remove' ||
          currentCommand === 'List') && (
          <Input
            placeholder="Path"
            width="auto"
            onChange={e => {
              setCommandData({ ...commandData, path: e.target.value });
            }}
          />
        )}

        <Button
          colorScheme="purple"
          size="lg"
          width="auto"
          onClick={submitCommand}
        >
          Enter Command
        </Button>
      </HStack>
    </>
  );
}
