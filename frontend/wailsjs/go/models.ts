export namespace models {
	
	export class Word {
	    id: number;
	    word: string;
	    phonetic: string;
	    pronunciation: string;
	    definition: string;
	    example: string;
	    translation: string;
	    imageUrl: string;
	    difficulty: number;
	    lastReviewed: number;
	    nextReview: number;
	    reviewCount: number;
	    easeFactor: number;
	    interval: number;
	    learned: boolean;
	    mastered: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Word(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.word = source["word"];
	        this.phonetic = source["phonetic"];
	        this.pronunciation = source["pronunciation"];
	        this.definition = source["definition"];
	        this.example = source["example"];
	        this.translation = source["translation"];
	        this.imageUrl = source["imageUrl"];
	        this.difficulty = source["difficulty"];
	        this.lastReviewed = source["lastReviewed"];
	        this.nextReview = source["nextReview"];
	        this.reviewCount = source["reviewCount"];
	        this.easeFactor = source["easeFactor"];
	        this.interval = source["interval"];
	        this.learned = source["learned"];
	        this.mastered = source["mastered"];
	    }
	}

}

