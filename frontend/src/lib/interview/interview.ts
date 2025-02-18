type Interview = {
    uuid: string,
    title: string,
    complete: boolean,
    timing: number,
    sections: Section[],
};

type Topic = {
    name: string
    grade: Grade
}

type Section = {
    uuid: string,
    name: string,
    complete: boolean,
    questions: Question[]
    position: number,
    grade: Grade,
}

type Question = {
    uuid: string,
    text: string,
    answer: string,
    grade: Grade
}

enum Grade {
    Junior = 'junior',
    Middle = 'middle',
    Senior = 'senior'
}

export {
    Grade
};

export type {
    Interview,
    Topic
}

