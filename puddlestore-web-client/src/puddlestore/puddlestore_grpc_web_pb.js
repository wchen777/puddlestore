/**
 * @fileoverview gRPC-Web generated client stub for pkg
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!


/* eslint-disable */
// @ts-nocheck



const grpc = {};
grpc.web = require('grpc-web');

const proto = {};
proto.pkg = require('./puddlestore_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?grpc.web.ClientOptions} options
 * @constructor
 * @struct
 * @final
 */
proto.pkg.PuddleStoreClient =
    function(hostname, credentials, options) {
  if (!options) options = {};
  options.format = 'text';

  /**
   * @private @const {!grpc.web.GrpcWebClientBase} The client
   */
  this.client_ = new grpc.web.GrpcWebClientBase(options);

  /**
   * @private @const {string} The hostname
   */
  this.hostname_ = hostname;

};


/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?grpc.web.ClientOptions} options
 * @constructor
 * @struct
 * @final
 */
proto.pkg.PuddleStorePromiseClient =
    function(hostname, credentials, options) {
  if (!options) options = {};
  options.format = 'text';

  /**
   * @private @const {!grpc.web.GrpcWebClientBase} The client
   */
  this.client_ = new grpc.web.GrpcWebClientBase(options);

  /**
   * @private @const {string} The hostname
   */
  this.hostname_ = hostname;

};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.pkg.Empty,
 *   !proto.pkg.ClientID>}
 */
const methodDescriptor_PuddleStore_ClientConnect = new grpc.web.MethodDescriptor(
  '/pkg.PuddleStore/ClientConnect',
  grpc.web.MethodType.UNARY,
  proto.pkg.Empty,
  proto.pkg.ClientID,
  /**
   * @param {!proto.pkg.Empty} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.pkg.ClientID.deserializeBinary
);


/**
 * @param {!proto.pkg.Empty} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.pkg.ClientID)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.pkg.ClientID>|undefined}
 *     The XHR Node Readable Stream
 */
proto.pkg.PuddleStoreClient.prototype.clientConnect =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/pkg.PuddleStore/ClientConnect',
      request,
      metadata || {},
      methodDescriptor_PuddleStore_ClientConnect,
      callback);
};


/**
 * @param {!proto.pkg.Empty} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.pkg.ClientID>}
 *     Promise that resolves to the response
 */
proto.pkg.PuddleStorePromiseClient.prototype.clientConnect =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/pkg.PuddleStore/ClientConnect',
      request,
      metadata || {},
      methodDescriptor_PuddleStore_ClientConnect);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.pkg.ClientID,
 *   !proto.pkg.Success>}
 */
const methodDescriptor_PuddleStore_ClientExit = new grpc.web.MethodDescriptor(
  '/pkg.PuddleStore/ClientExit',
  grpc.web.MethodType.UNARY,
  proto.pkg.ClientID,
  proto.pkg.Success,
  /**
   * @param {!proto.pkg.ClientID} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.pkg.Success.deserializeBinary
);


/**
 * @param {!proto.pkg.ClientID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.pkg.Success)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.pkg.Success>|undefined}
 *     The XHR Node Readable Stream
 */
proto.pkg.PuddleStoreClient.prototype.clientExit =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/pkg.PuddleStore/ClientExit',
      request,
      metadata || {},
      methodDescriptor_PuddleStore_ClientExit,
      callback);
};


/**
 * @param {!proto.pkg.ClientID} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.pkg.Success>}
 *     Promise that resolves to the response
 */
proto.pkg.PuddleStorePromiseClient.prototype.clientExit =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/pkg.PuddleStore/ClientExit',
      request,
      metadata || {},
      methodDescriptor_PuddleStore_ClientExit);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.pkg.OpenMessage,
 *   !proto.pkg.OpenResponse>}
 */
const methodDescriptor_PuddleStore_ClientOpen = new grpc.web.MethodDescriptor(
  '/pkg.PuddleStore/ClientOpen',
  grpc.web.MethodType.UNARY,
  proto.pkg.OpenMessage,
  proto.pkg.OpenResponse,
  /**
   * @param {!proto.pkg.OpenMessage} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.pkg.OpenResponse.deserializeBinary
);


/**
 * @param {!proto.pkg.OpenMessage} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.pkg.OpenResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.pkg.OpenResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.pkg.PuddleStoreClient.prototype.clientOpen =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/pkg.PuddleStore/ClientOpen',
      request,
      metadata || {},
      methodDescriptor_PuddleStore_ClientOpen,
      callback);
};


/**
 * @param {!proto.pkg.OpenMessage} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.pkg.OpenResponse>}
 *     Promise that resolves to the response
 */
proto.pkg.PuddleStorePromiseClient.prototype.clientOpen =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/pkg.PuddleStore/ClientOpen',
      request,
      metadata || {},
      methodDescriptor_PuddleStore_ClientOpen);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.pkg.CloseMessage,
 *   !proto.pkg.Success>}
 */
const methodDescriptor_PuddleStore_ClientClose = new grpc.web.MethodDescriptor(
  '/pkg.PuddleStore/ClientClose',
  grpc.web.MethodType.UNARY,
  proto.pkg.CloseMessage,
  proto.pkg.Success,
  /**
   * @param {!proto.pkg.CloseMessage} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.pkg.Success.deserializeBinary
);


/**
 * @param {!proto.pkg.CloseMessage} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.pkg.Success)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.pkg.Success>|undefined}
 *     The XHR Node Readable Stream
 */
proto.pkg.PuddleStoreClient.prototype.clientClose =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/pkg.PuddleStore/ClientClose',
      request,
      metadata || {},
      methodDescriptor_PuddleStore_ClientClose,
      callback);
};


/**
 * @param {!proto.pkg.CloseMessage} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.pkg.Success>}
 *     Promise that resolves to the response
 */
proto.pkg.PuddleStorePromiseClient.prototype.clientClose =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/pkg.PuddleStore/ClientClose',
      request,
      metadata || {},
      methodDescriptor_PuddleStore_ClientClose);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.pkg.WriteMessage,
 *   !proto.pkg.Success>}
 */
const methodDescriptor_PuddleStore_ClientWrite = new grpc.web.MethodDescriptor(
  '/pkg.PuddleStore/ClientWrite',
  grpc.web.MethodType.UNARY,
  proto.pkg.WriteMessage,
  proto.pkg.Success,
  /**
   * @param {!proto.pkg.WriteMessage} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.pkg.Success.deserializeBinary
);


/**
 * @param {!proto.pkg.WriteMessage} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.pkg.Success)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.pkg.Success>|undefined}
 *     The XHR Node Readable Stream
 */
proto.pkg.PuddleStoreClient.prototype.clientWrite =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/pkg.PuddleStore/ClientWrite',
      request,
      metadata || {},
      methodDescriptor_PuddleStore_ClientWrite,
      callback);
};


/**
 * @param {!proto.pkg.WriteMessage} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.pkg.Success>}
 *     Promise that resolves to the response
 */
proto.pkg.PuddleStorePromiseClient.prototype.clientWrite =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/pkg.PuddleStore/ClientWrite',
      request,
      metadata || {},
      methodDescriptor_PuddleStore_ClientWrite);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.pkg.ReadMessage,
 *   !proto.pkg.ReadResponse>}
 */
const methodDescriptor_PuddleStore_ClientRead = new grpc.web.MethodDescriptor(
  '/pkg.PuddleStore/ClientRead',
  grpc.web.MethodType.UNARY,
  proto.pkg.ReadMessage,
  proto.pkg.ReadResponse,
  /**
   * @param {!proto.pkg.ReadMessage} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.pkg.ReadResponse.deserializeBinary
);


/**
 * @param {!proto.pkg.ReadMessage} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.pkg.ReadResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.pkg.ReadResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.pkg.PuddleStoreClient.prototype.clientRead =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/pkg.PuddleStore/ClientRead',
      request,
      metadata || {},
      methodDescriptor_PuddleStore_ClientRead,
      callback);
};


/**
 * @param {!proto.pkg.ReadMessage} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.pkg.ReadResponse>}
 *     Promise that resolves to the response
 */
proto.pkg.PuddleStorePromiseClient.prototype.clientRead =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/pkg.PuddleStore/ClientRead',
      request,
      metadata || {},
      methodDescriptor_PuddleStore_ClientRead);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.pkg.MkdirMessage,
 *   !proto.pkg.Success>}
 */
const methodDescriptor_PuddleStore_ClientMkdir = new grpc.web.MethodDescriptor(
  '/pkg.PuddleStore/ClientMkdir',
  grpc.web.MethodType.UNARY,
  proto.pkg.MkdirMessage,
  proto.pkg.Success,
  /**
   * @param {!proto.pkg.MkdirMessage} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.pkg.Success.deserializeBinary
);


/**
 * @param {!proto.pkg.MkdirMessage} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.pkg.Success)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.pkg.Success>|undefined}
 *     The XHR Node Readable Stream
 */
proto.pkg.PuddleStoreClient.prototype.clientMkdir =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/pkg.PuddleStore/ClientMkdir',
      request,
      metadata || {},
      methodDescriptor_PuddleStore_ClientMkdir,
      callback);
};


/**
 * @param {!proto.pkg.MkdirMessage} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.pkg.Success>}
 *     Promise that resolves to the response
 */
proto.pkg.PuddleStorePromiseClient.prototype.clientMkdir =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/pkg.PuddleStore/ClientMkdir',
      request,
      metadata || {},
      methodDescriptor_PuddleStore_ClientMkdir);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.pkg.RemoveMessage,
 *   !proto.pkg.Success>}
 */
const methodDescriptor_PuddleStore_ClientRemove = new grpc.web.MethodDescriptor(
  '/pkg.PuddleStore/ClientRemove',
  grpc.web.MethodType.UNARY,
  proto.pkg.RemoveMessage,
  proto.pkg.Success,
  /**
   * @param {!proto.pkg.RemoveMessage} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.pkg.Success.deserializeBinary
);


/**
 * @param {!proto.pkg.RemoveMessage} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.pkg.Success)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.pkg.Success>|undefined}
 *     The XHR Node Readable Stream
 */
proto.pkg.PuddleStoreClient.prototype.clientRemove =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/pkg.PuddleStore/ClientRemove',
      request,
      metadata || {},
      methodDescriptor_PuddleStore_ClientRemove,
      callback);
};


/**
 * @param {!proto.pkg.RemoveMessage} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.pkg.Success>}
 *     Promise that resolves to the response
 */
proto.pkg.PuddleStorePromiseClient.prototype.clientRemove =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/pkg.PuddleStore/ClientRemove',
      request,
      metadata || {},
      methodDescriptor_PuddleStore_ClientRemove);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.pkg.ListMessage,
 *   !proto.pkg.ListResponse>}
 */
const methodDescriptor_PuddleStore_ClientList = new grpc.web.MethodDescriptor(
  '/pkg.PuddleStore/ClientList',
  grpc.web.MethodType.UNARY,
  proto.pkg.ListMessage,
  proto.pkg.ListResponse,
  /**
   * @param {!proto.pkg.ListMessage} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.pkg.ListResponse.deserializeBinary
);


/**
 * @param {!proto.pkg.ListMessage} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.pkg.ListResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.pkg.ListResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.pkg.PuddleStoreClient.prototype.clientList =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/pkg.PuddleStore/ClientList',
      request,
      metadata || {},
      methodDescriptor_PuddleStore_ClientList,
      callback);
};


/**
 * @param {!proto.pkg.ListMessage} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.pkg.ListResponse>}
 *     Promise that resolves to the response
 */
proto.pkg.PuddleStorePromiseClient.prototype.clientList =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/pkg.PuddleStore/ClientList',
      request,
      metadata || {},
      methodDescriptor_PuddleStore_ClientList);
};


module.exports = proto.pkg;

