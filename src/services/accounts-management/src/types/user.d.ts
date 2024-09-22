export type UserAccount = {
    username: string,
    credentials: { password?: string, }
    userid: string,
    active: boolean,
    record: {
        createdOn: Date
    }
}