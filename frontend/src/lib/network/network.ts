import {Interview, Question, Topic} from "@/lib/interview/interview.ts";
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


type AnswerQuestionResponse = {
    updatedQuestion: Question|null,
    nextSection: boolean,
}

class Network {

    private serviceUrl: string;

    constructor(serviceUrl: string, isDev: boolean) {
        this.serviceUrl = serviceUrl;
        if (isDev) {
            console.log('network: using dev mode');
        }
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

    async answerQuestion(questionId: string, answer: string): Promise<AnswerQuestionResponse> {
        const {data, status} = await axios.post<Question>(
            this.serviceUrl + '/api/question',
            {
                answer: answer,
                question_uuid: questionId,
            },
            {
                headers: {
                    Accept: 'application/json',
                    'Content-Type': 'application/json',
                },
                withCredentials: true,
            }
        )

        if (status == 205)
        {
            location.reload();
        }



        if (status == 220)
        {
            return {
                updatedQuestion: data,
                nextSection: true,
            }
        } else if (status == 200) {
            return {
                updatedQuestion: data,
                nextSection: false,
            }
        }

        throw new Error("Invalid code while get interview list: expeted 200/220/205, got " + status);
    }

    async getNextQuestion(interviewId: string): Promise<Question> {
        const {data, status} = await axios.get<Question>(
            this.serviceUrl + '/api/question/next/' + interviewId,
            {
                headers: {
                    Accept: 'application/json',
                    'Content-Type': 'application/json',
                },
                withCredentials: true,
            }
        )

        if (status == 205)
        {
            location.reload();
        }

        if (status == 200) {
            return data;
        }

        throw new Error("Invalid code while get next question: expected 200, got " + status);
    }

    async getNextSectionQuestion(interviewId: string): Promise<Question> {
        const {data, status} = await axios.get<Question>(
            this.serviceUrl + '/api/question/change/' + interviewId,
            {
                headers: {
                    Accept: 'application/json',
                    'Content-Type': 'application/json',
                },
                withCredentials: true,
            }
        )

        if (status == 205)
        {
            location.reload();
        }

        if (status == 200) {
            return data;
        }

        throw new Error("Invalid code while get next section question: expected 200, got " + status);
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

    async createInterviewFeedback(interviewUUID: string): Promise<string> {
        const {data, status} = await axios.post<string>(
            this.serviceUrl + '/interview/feedback/' + interviewUUID,
            interviewUUID,
            {
                withCredentials: true,
            }
        )

        if (status != 200 && status != 201) {
            throw new Error("Invalid code while preview message: expeted 200/201, got " + status);
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