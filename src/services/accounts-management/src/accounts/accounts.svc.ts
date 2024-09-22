import { UserAccount } from "../types/user"

export const createAccount = async (userAccount: UserAccount): Promise<boolean> => {

    return false
}

export const changeAccountStatus = async (userid: string, enable: boolean): Promise<boolean> => {

    return false
}

export const getAccount = async (userid: string): Promise<UserAccount> => {

    return {} as UserAccount;
}