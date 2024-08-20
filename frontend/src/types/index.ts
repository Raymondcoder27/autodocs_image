
import type { Themes } from "@/constants"
export type TTheme = (typeof Themes)[number]


export interface ApiResponse<T> {
    message: string
    status:number
    error?:string
    data?: T
}

export interface ApiErrorResponse {
    data:object
    message:string
    error:string
    status:number
}
