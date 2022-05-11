import React from 'react';
import { Text, UnorderedList, ListItem, Box } from '@chakra-ui/react';
import './Main.css';

export default function OpenData({ openFDs }) {
  return (
    <div className="left-justify">
      <Text fontSize="xl" mb={3}>
        Open File Descriptors
      </Text>

      <Box
        overflowY={'scroll'}
        className="border"
        minHeight={70}
        width={'30vw'}
        maxHeight={'40vh'}
      >
        {openFDs.length === 0 ? (
          <Text p={2}>No open file data</Text>
        ) : (
          <UnorderedList>
            {openFDs.map(fd => {
              return (
                <ListItem>
                  {fd.fd}: {fd.filepath}
                </ListItem>
              );
            })}
          </UnorderedList>
        )}
      </Box>
    </div>
  );
}
