/* Package limiters provides different Job execution limits for Workers.

It does define the Limiter interface that should be implemented on every
limiter as the WOrker would loop for every configured limiter to aquire or run
the job.

Initially the limiters package have been created to avoid hitting throttling
limits on the Twitter's API.

*/
package limiters
