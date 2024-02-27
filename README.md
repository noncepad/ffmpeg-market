# FFmpeg-market

this code does this

for details go to noncepad.com

# Build


# Run



FFmpegConverter

How Solpipe empowered me to harvest and transform my internal tool into a marketable serivce.

Harvesting the Gold from Within 

Introduction:

In today's rapidly evolving business landscape, organizations strive to optimize internal processes for greater efficiency and productivity. Like many, i was given a task from my boss, saw the benefits in automating said task and develpoed a tool to do so. Seeing its potential, I was able to utilize Solpipe to easily transform my tool into valuable assets for the broader market.

The Challenge:

It all began with what seemed to be a simple task: rendering blender files (digital content creation software) into gifs for our website's homepage, using FFmpeg. I, previously unfamiliar with FFmpeg, immidietly noticed some drawbacks with ffmpeg commands presenting complexities and resource-intensive demands. Rendering these videos took an excessive amount of compute power. And repeating this process for all the videos I was tasked to render would exhauste my resources. This made me realize that outsourcing this process, especially for my larger blender files, would be a much better option.

The Idea Takes Shape:
Faced with the prospect of repeating this process countless times, the seed of an idea sprouted. Why not build a tool to automate file conversion? And then, why not just build it for myself, but for others facing similar challenges? 

Overcoming Hurdles:
I outlined what hurdles I would need to overcome to bring this idea a fruition:
1. I would need to see if I could use FFmpeg commands in my coding language, Go. 
-----hurdle1clip1 from graphics
2. How can I daemonize this process and build a distributed system?
3. How can I distribute my tool and make money?

Researching ffmpeg integration in Go and exploring daemonization options marked the initial steps. I was very easily able to write a go test and confirm that i would be able to create an FFmpeg wrapper in Go. For daemonization i decided to implement a manager-to-worker model (link to page about this model here) and craft a protobuf for a grpc server. 

This way I could allow for scalability to handle increased workloads concurrently, adding more workers as needed. I opted to run my server over gRPC, leveraging the straightforward and intuitive data format of Protobufs to serialize my structured data. This strategic choice aligns with the potential to monetize my software using Solpipe, which offers seamless integration and accessibility to a broader market.



(something in this story has to link hurdle 2 and 3. Using protobuf and grpc would be the easiest way to monetize, ie. using solpipe. 
)

Daemonizing with gRPC and Protocol Buffers:


Monetization Strategy:

In navigating the monetization landscape for my software, I initially explored the avenues of Stripe and AWS, seeking to leverage their robust platforms for distributing my API. However, I encountered complexities and barriers inherent in these services, ranging from intricate setup processes to scalability challenges. 
* Both Stripe and AWS require intricate setup processes, involving account registration, configuration, and integration into existing systems.
* Both Stripe and AWS have associated costs, including setup fees, transaction fees, and ongoing service charges.
* Deciding the pricing structures proved arduous and time consuming.

On platforms like Stripe and AWS, pricing is decided by a human and presents a multifaceted challenge when entering a market. With Stripe, intricacies arise in determining transaction fees, subscription tiers, and discount structures, all while balancing competitiveness and profitability. Similarly, AWS offers a plethora of services with varied pricing models requiring careful consideration of usage patterns and resource requirements. Navigating these complexities demands a deep understanding of both the product offerings and the market landscape, as well as ongoing analysis to optimize pricing strategies and maintain competitiveness in a dynamic environment.

Recognizing the need for a more accessible solution, I turned to Solpipe.With Solpipe, instead of directly deciding on pricing, I set constraints on the capacity I could offer, allowing the price to be determined by the demand and the resources available. This approach was particularly beneficial in the context of platforms like Stripe and AWS, where the intricate pricing structures can be daunting to navigate. By setting clear capacity limitations, easily managed by my worker-pool, I streamlined the process and let market forces dictate the pricing, ensuring a fair and competitive model that aligns with both the value provided and the resources expended.

Monetizing Solution: 

Solpipes intuitive interface and streamlined processes, automates pricing and offers a refreshing contrast. Solpipe simplifies the deployment of an API. I was able to install the necessary executables, follow straightforward instructions, and swiftly post my offerings to the marketplace. By running a pipeline, I provided others with seamless access to my API, all without the complexities and hurdles encountered with traditional platforms. With Solpipe, market entry became swift; the path to monetization became clear, efficient, and hassle-free.