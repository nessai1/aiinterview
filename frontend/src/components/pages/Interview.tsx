import React from "react";
import {Link, useParams} from "react-router-dom";
import Chat from "@/components/internal/chat/Chat.tsx";
import {ChevronLeft} from 'lucide-react';
import { Button } from "@/components/ui/button";



const Interview: React.FC = () => {

    const { interviewId } = useParams();
    if (interviewId === undefined) {
        window.location.href = '/';
        return;
    }

    return (
        <>
            <Link to={"/"} style={{ textDecoration: "none", color: "inherit"}}>
                <Button variant="outline" size="icon" className="mb-2">
                    <ChevronLeft />
                </Button>
            </Link>
            <Chat interviewId={interviewId} />
        </>
    );
}

export default Interview;