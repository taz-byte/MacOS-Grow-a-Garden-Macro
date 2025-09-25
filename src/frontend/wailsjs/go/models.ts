export namespace settingsmanager {
	
	export class Settings {
	    buy_carrot: boolean;
	    buy_strawberry: boolean;
	    buy_blueberry: boolean;
	    buy_orange_tulip: boolean;
	    buy_tomato: boolean;
	    buy_corn: boolean;
	    buy_daffodil: boolean;
	    buy_watermelon: boolean;
	    buy_pumpkin: boolean;
	    buy_apple: boolean;
	    buy_bamboo: boolean;
	    buy_coconut: boolean;
	    buy_cactus: boolean;
	    buy_dragon_fruit: boolean;
	    buy_mango: boolean;
	    buy_grape: boolean;
	    buy_mushroom: boolean;
	    buy_pepper: boolean;
	    buy_cacao: boolean;
	    buy_beanstalk: boolean;
	    buy_ember_lily: boolean;
	    buy_sugar_apple: boolean;
	    buy_burning_bud: boolean;
	    buy_giant_pinecone: boolean;
	    buy_elder_strawberry: boolean;
	    buy_romanesco: boolean;
	    enable_discord_webhook: boolean;
	    discord_webhook_url: string;
	
	    static createFrom(source: any = {}) {
	        return new Settings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.buy_carrot = source["buy_carrot"];
	        this.buy_strawberry = source["buy_strawberry"];
	        this.buy_blueberry = source["buy_blueberry"];
	        this.buy_orange_tulip = source["buy_orange_tulip"];
	        this.buy_tomato = source["buy_tomato"];
	        this.buy_corn = source["buy_corn"];
	        this.buy_daffodil = source["buy_daffodil"];
	        this.buy_watermelon = source["buy_watermelon"];
	        this.buy_pumpkin = source["buy_pumpkin"];
	        this.buy_apple = source["buy_apple"];
	        this.buy_bamboo = source["buy_bamboo"];
	        this.buy_coconut = source["buy_coconut"];
	        this.buy_cactus = source["buy_cactus"];
	        this.buy_dragon_fruit = source["buy_dragon_fruit"];
	        this.buy_mango = source["buy_mango"];
	        this.buy_grape = source["buy_grape"];
	        this.buy_mushroom = source["buy_mushroom"];
	        this.buy_pepper = source["buy_pepper"];
	        this.buy_cacao = source["buy_cacao"];
	        this.buy_beanstalk = source["buy_beanstalk"];
	        this.buy_ember_lily = source["buy_ember_lily"];
	        this.buy_sugar_apple = source["buy_sugar_apple"];
	        this.buy_burning_bud = source["buy_burning_bud"];
	        this.buy_giant_pinecone = source["buy_giant_pinecone"];
	        this.buy_elder_strawberry = source["buy_elder_strawberry"];
	        this.buy_romanesco = source["buy_romanesco"];
	        this.enable_discord_webhook = source["enable_discord_webhook"];
	        this.discord_webhook_url = source["discord_webhook_url"];
	    }
	}

}

