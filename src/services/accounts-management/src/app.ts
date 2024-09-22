// import express, { Express } from "express";
import express, { Express, Request, Response } from "express";
import accountController from "./accounts/accounts.ctr"

export const app = async () => {
    const app: Express = express();
    app.use('/account', accountController)
    app.listen(4001)
}