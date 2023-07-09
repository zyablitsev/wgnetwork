'use strict'

class RPCError extends Error {
    constructor(code, message) {
        super(`${code}: ${message}`);
        this.code = code;
        this.message = message;
    }

    get name() { return "RPCError"; }
}

class ValidationError extends Error {
    constructor(field, message) {
        super(`${field}: ${message}`);
        this.field = field;
        this.message = message;
    }

    get name() { return "ValidationError"; }
}

class RPCApi {
  constructor(url) {
    this.url = new URL(url);
  }

  async Fetch(method, params, session) {
    const body = {
      'method': method,
      'params': params
    }

    let headers = new Headers({
      'accept': 'application/json',
      'content-type': 'application/json'
    });

    if (typeof session !== 'undefined') {
      headers.append('x-session', session);
    }

    const cfg = {
      'method': 'POST',
      'headers': headers,
      'body': JSON.stringify(body)
    };

    let request = new Request(this.url, cfg);

    try {
      const response = await fetch(request);
      const data = await response.json();

      if (response.status !== 200) {
        let error = new RPCError(400, response.statusText, undefined);
        return Promise.reject(error);
      }

      let type = response.headers.get("content-type");
      if (type !== "application/json") {
        let error = new RPCError(400, `expected json, got ${type}`, undefined);
        return Promise.reject(error);
      }

      if (typeof data.error !== 'undefined') {
        if (data.error.code === 401) {
          let verr = new ValidationError(
            data.error.data.field,
            data.error.data.message);
          return Promise.reject(verr);
        } else {
          let rpcerr = new RPCError(data.error.code, data.error.message);
          return Promise.reject(rpcerr);
        }
      }

      return Promise.resolve(data.result);
    } catch (err) {
      let msg = 'unexpected error'

      if (err instanceof NetworkError) {
        msg = 'network error'
      } else if (err instanceof AbortError) {
        msg = 'request canceled'
      }

      let error = new RPCError(400, msg, undefined);

      return Promise.reject(error);
    }
  }
}

export default RPCApi;
export { RPCError, ValidationError };
