import React, {useEffect} from "react";
import {CountdownTimer} from "@/components/internal/chat/CountdownTimer.tsx";
import SectionBlock from "@/components/internal/chat/SectionBlock.tsx";
import {Interview} from "@/lib/interview/interview.ts";
import {AxiosError} from "axios";
import {ToastAction} from "@/components/ui/toast.tsx";


type TProps = {
    interviewId: string
}

const Chat: React.FC<TProps> = (props: TProps) => {


    const retry = () => {
        network.getInterviewList().then((interviews: Interview[]) => {
            console.log(interviews);
            setItems(interviews);
        }).catch((err: AxiosError) => {
            toast({
                title: 'Упс! Список интервью не загрузился ;(',
                description: `Ошибка сети: [${err.code}] ${err.message}`,
                action: (
                    <ToastAction altText="Goto schedule to undo" onClick={retry}>Повторить</ToastAction>
                ),
                variant: "destructive",
            });
        });
    };

    useEffect(() => {
        retry();
    }, []);

    return (
        <>
            {/*<CountdownTimer seconds={180}/>*/}
            {/*Interview: {props.interviewId}  Interview: {props.interviewId}  Interview: {props.interviewId}*/}
            <SectionBlock section={{}} />
        </>
    );
}

export default Chat;