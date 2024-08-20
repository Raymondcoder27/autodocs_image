export interface Doc{
    id:string
    documentName:string
    templateId:string
    data:string
    refNumber: string
    created_at: string
}

export interface GenerationRequest {
    refNumber:string
    description:string
    data:string
}