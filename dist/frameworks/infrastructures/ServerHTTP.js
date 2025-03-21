"use strict";
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (this && this.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};
var __param = (this && this.__param) || function (paramIndex, decorator) {
    return function (target, key) { decorator(target, key, paramIndex); }
};
Object.defineProperty(exports, "__esModule", { value: true });
const common_1 = require("@nestjs/common");
let ServerHTTP = class ServerHTTP {
    constructor(httpServer) {
        this.httpServer = httpServer;
    }
    async get(url, options) {
        return this.httpServer(url, { ...options, method: "GET" });
    }
    async post(url, body, options) {
        return this.httpServer(url, {
            method: "POST",
            body: JSON.stringify(body),
            headers: { "Content-Type": "application/json", ...options?.headers },
            ...options
        });
    }
    async put(url, body, options) {
        return this.httpServer(url, {
            method: "PUT",
            body: JSON.stringify(body),
            headers: { "Content-Type": "application/json", ...options?.headers },
            ...options
        });
    }
    async delete(url, options) {
        return this.httpServer(url, { ...options, method: "DELETE" });
    }
};
ServerHTTP = __decorate([
    (0, common_1.Injectable)(),
    __param(0, (0, common_1.Inject)("IHttpServer")),
    __metadata("design:paramtypes", [Function])
], ServerHTTP);
exports.default = ServerHTTP;
//# sourceMappingURL=ServerHTTP.js.map