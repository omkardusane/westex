import { Request, Response, Router } from "express";
import { createAccount } from "./accounts.svc"

const controller: Router = Router();
export default controller;

controller.get('/one', async (req: Request, res: Response) => {

})

controller.post('/one', async (req: Request, res: Response) => {

})

controller.get('/all', async (req: Request, res: Response) => {

})