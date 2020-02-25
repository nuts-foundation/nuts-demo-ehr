import io from 'socket.io-client';

const consent    = io('/consent');
const accessLogs = io('/accessLogs');

export default {
  consent:    () => consent,
  accessLogs: () => accessLogs
};
