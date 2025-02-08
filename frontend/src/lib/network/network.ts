import {Interview} from "@/lib/interview/interview.ts";
import axios from 'axios';

type GetInterviewListResponse = Interview[]

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
}

export {Network};