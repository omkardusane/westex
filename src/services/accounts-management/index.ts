import { app } from "./src/app"
import { logger } from "../../shared/logger"

app().then(() => {
    logger.info('Accounts service online')
}).catch((error) => {
    logger.error('Accounts service: error', error)
});