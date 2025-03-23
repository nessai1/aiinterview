import {Interview, Topic} from "@/lib/interview/interview.ts";
import axios from 'axios';

type GetInterviewListResponse = Interview[]
type GetInterviewResponse = Interview

type CreateInterviewRequest = {
    title: string,
    timing: number,
    topics: Topic[]
}

type CreateInterviewResponse = {
    uuid: string
}

class Network {

    private serviceUrl: string
    private isDev: boolean;

    constructor(serviceUrl: string, isDev: boolean) {
        this.serviceUrl = serviceUrl;
        this.isDev = isDev;
    }

    async getInterviewList(): Promise<Interview[]> {
       const {data, status} = await axios.get<GetInterviewListResponse>(
           this.serviceUrl + '/api/interview/list',
           {
               headers: {
                   Accept: 'application/json',
               },
               withCredentials: true,
           }
       )

        if (status !== 200) {
            throw new Error("Invalid code while get interview list: expeted 200, got " + status);
        }

        return data;
    }

    async loadInterview(interviewId: string): Promise<Interview> {
        const {data, status} = await axios.get<GetInterviewResponse>(
            this.serviceUrl + '/api/interview/' + interviewId,
            {
                headers: {
                    Accept: 'application/json',
                },
                withCredentials: true,
            }
        )

        if (status !== 200) {
            throw new Error("Invalid code while get interview list: expeted 200, got " + status);
        }

        return data;
    }

    async createInterview(interview: CreateInterviewRequest): Promise<CreateInterviewResponse> {
        const {data, status} = await axios.post<CreateInterviewResponse>(
            this.serviceUrl + '/api/interview',
            interview,
            {
                headers: {
                    Accept: 'application/json',
                    'Content-Type': 'application/json',
                },
                withCredentials: true,
            }
        )

        if (status !== 200) {
            throw new Error("Invalid code while create interview: expeted 200, got " + status);
        }

        return data;
    }

    async previewMessage(message: string): Promise<string> {
        const {data, status} = await axios.post<string>(
            this.serviceUrl + '/api/preview',
            message,
            {
                withCredentials: true,
            }
        )

        if (status !== 200) {
            throw new Error("Invalid code while preview message: expeted 200, got " + status);
        }

        return data;
    }
}

export {Network};