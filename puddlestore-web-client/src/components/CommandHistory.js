import React from 'react';
import {
  Box,
  Text,
  VStack,
  HStack,
  Divider,
  Stat,
  StatLabel,
  StatHelpText,
  StatNumber,
  Badge,
} from '@chakra-ui/react';
import './Main.css';

export default function CommandHistory({ commandHistory }) {
  return (
    <div className="left-justify">
      <Text fontSize="xl" mb={3}>
        Command History
      </Text>

      <Box
        className="border"
        width={'50vw'}
        padding={5}
        minHeight={70}
        maxHeight={'50vh'}
        overflowY={'scroll'}
      >
        {commandHistory.map(c => {
          return (
            <>
              <Stat>
                <StatLabel>
                  <b>{c.command}</b>{' '}
                  {c.result === 'success' && (
                    <Badge colorScheme="green">Success</Badge>
                  )}
                </StatLabel>
                <StatHelpText>{c.argString}</StatHelpText>
                {c.result === 'success' && <StatLabel>{c.result}</StatLabel>}
              </Stat>
              <Divider orientation="horizontal" />
            </>
          );
        })}
      </Box>
    </div>
  );
}
