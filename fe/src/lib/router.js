'use strict'

// Router class.
class Router {
  constructor() {
    this.Patterns = {};
    this.Static = {};
    this.ParametrizedTree = {};
    this.ParametrizedCache = {};
    this.NameRouteMap = {};
  }

  Add(pattern, name, handler) {
    let route, depth, parametrizedNode;
    if ((typeof pattern === 'undefined' || pattern.charAt(0) !== '/') ||
        (typeof name === 'undefined' || name == '')) {
      return;
    }
    route = new Route(pattern, handler, name);
    this.add(
      route,
      this.Patterns,
      this.Static,
      this.ParametrizedTree,
    );
  }

  Match(path) {
    return this.match(
      path,
      this.ParametrizedCache,
      this.Static,
      this.ParametrizedTree,
    );
  }

  ReverseURI(name, params, query) {
    let route = this.NameRouteMap[name];
    if (typeof route === 'undefined') return;
    query = (query instanceof URLSearchParams) ? '?' + query.toString() : '';
    if (!route.isParametrized) return route.pattern + query;
    let parts = [''];
    parts.push(...route.parametrized.patternParts);
    for (let i = 0; i < route.parametrized.paramsLength; i++) {
      let value = params[route.parametrized.params[i].name];
      if (typeof value === 'undefined') return;
      parts[route.parametrized.params[i].position+1] = value;
    }
    return parts.join('/') + query;
  }

  add(route, patterns, staticMap, parametrizedNode) {
    let depth;
    if (typeof this.NameRouteMap[route.name] !== 'undefined' ||
        typeof patterns[route.pattern] !== 'undefined') {
      return;
    }
    this.NameRouteMap[route.name] = route;
    patterns[route.pattern] = true;
    if (!route.isParametrized) {
      staticMap[route.pattern] = route;
      return;
    }
    // build parametrized tree
    depth = route.parametrized.patternPartsLength;
    for (let i = 0; i < route.parametrized.patternPartsLength; i++) {
      depth--;
      if (typeof parametrizedNode[route.parametrized.patternParts[i]] === 'undefined') {
        parametrizedNode[route.parametrized.patternParts[i]] = {
          route: (depth === 0) ? route : undefined,
          patternPart: route.parametrized.patternParts[i],
          child: undefined,
          depth: depth};
      } else {
        if (parametrizedNode[route.parametrized.patternParts[i]].depth < depth) {
          parametrizedNode[route.parametrized.patternParts[i]].depth = depth;
        }
        if (depth === 0) {
          parametrizedNode[route.parametrized.patternParts[i]].route = route;
          break;
        }
      }
      if (depth > 0 && typeof parametrizedNode[route.parametrized.patternParts[i]].child === 'undefined') {
        parametrizedNode[route.parametrized.patternParts[i]].child = {};
      }
      parametrizedNode = parametrizedNode[route.parametrized.patternParts[i]].child;
    }
  }

  match(path, parametrizedCache, staticMap, parametrizedNode) {
    let route, queryPartIdx, href, query, pathParts, pathPartsLength, found; //, params;
    queryPartIdx = path.indexOf('?');
    query = '';
    href = path;
    if (queryPartIdx > -1) {
      query = path.substring(queryPartIdx+1);
      href = path.substring(0, queryPartIdx);
    }
    pathParts = href.split('/').slice(1);
    pathPartsLength = pathParts.length;
    if (typeof staticMap[href] !== 'undefined') { // check static
      return {
        handler: staticMap[href].handler,
        location: {
          pathname: href,
          search: query,
          name: staticMap[href].name}
      }
    }
    if (typeof parametrizedCache[href] !== 'undefined') { // check cache
      return {
        handler: parametrizedCache[href].route.handler,
        location: {
          pathname: href,
          search: query,
          name: parametrizedCache[href].route.name},
        params: parametrizedCache[href].params,
      }
    }
    for (let i = 0, childFlag = pathPartsLength-1; i < pathPartsLength; i++) {
      found = parametrizedNode[':'] || parametrizedNode[pathParts[i]];
      if (typeof found === 'undefined') break;
      if (i === childFlag) {
        route = found.route;
      }
      parametrizedNode = found.child
    }
    if (typeof route === 'undefined') return;

    // process params
    let params = {};
    for (let i = 0; i < route.parametrized.paramsLength; i++) {
      params[route.parametrized.params[i].name] = pathParts[route.parametrized.params[i].position];
    }

    if (typeof parametrizedCache[href] === 'undefined') {
      parametrizedCache[href] = {route: route, params: params};
    }

    return {
      handler: route.handler,
      location: {
        pathname: href,
        search: query,
        name: route.name},
      params: params,
    }
  }
}

// Route class.
class Route {
  constructor(pattern, handler, name) {
    if (typeof pattern === 'undefined') return;
    this.pattern = pattern;
    this.handler = handler;
    this.name = name;
    this.isParametrized = this.pattern.indexOf(':') < 0 ? false : true;
    if (this.isParametrized) this.parseParametrizedPattern();
  }

  parseParametrizedPattern() {
    if (!this.isParametrized) return;
    let parts = this.pattern.split('/').slice(1);
    this.parametrized = {
      patternParts: parts,
      patternPartsLength: parts.length,
      params: [],
      paramsMap: {},
      paramsLength: 0,
    }
    for (let i = 0; i < this.parametrized.patternPartsLength; i++) {
      if (this.parametrized.patternParts[i].charAt(0) !== ':') continue;
      let parameter = this.parametrized.patternParts[i].slice(1);
      if (typeof this.parametrized.paramsMap[parameter] !== 'undefined') continue;
      this.parametrized.paramsMap[parameter] = {
        position: i, name: parameter};
      this.parametrized.params.push(this.parametrized.paramsMap[parameter])
      this.parametrized.patternParts[i] = ':';
    }
    this.parametrized.paramsLength = this.parametrized.params.length;
    this.pattern = [''];
    this.pattern.push(...this.parametrized.patternParts);
    this.pattern = this.pattern.join('/');
  }
}

export default Router;
