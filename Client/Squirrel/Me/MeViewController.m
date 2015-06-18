//
//  MeViewController.m
//  Squirrel
//
//  Created by JamesMao on 4/11/15.
//  Copyright (c) 2015 JamesMao. All rights reserved.
//

#import "MeViewController.h"
#import "Define.h"
#import "SettingViewController.h"
#import "NotificationPublishViewController.h"
#import "AdminPublishTableViewCell.h"
#import "AFNetworking.h"
#import "MeViewData.h"

@interface MeViewController ()
{
    IBOutlet UITableView* classesTableView_;
}

@end

@implementation MeViewController

- (void)viewDidLoad
{
    [super viewDidLoad];
    
    self.navigationItem.rightBarButtonItem = [[UIBarButtonItem alloc] initWithTitle:@"Setting"
                                                                              style:UIBarButtonItemStyleDone
                                                                             target:self
                                                                             action:@selector(onSetting)];
}

- (void)onSetting
{
    [self.navigationController pushViewController:[SettingViewController singleton] animated:YES];
}

- (void)didReceiveMemoryWarning
{
    [super didReceiveMemoryWarning];
    // Dispose of any resources that can be recreated.
}

- (void)viewWillAppear:(BOOL)animated
{
    dispatch_async(BackGroundQueue, ^{
        NSString* urlString = [NSString stringWithFormat:@"%@%@%@", SERVER_IP, URL_PART_ONE_USER_CLASS, [MeViewData singleton].userUniqueName];
        NSURL* url = [NSURL URLWithString: urlString];
        NSData* data = [NSData dataWithContentsOfURL: url];
        [self performSelectorOnMainThread:@selector(fetchedData:) withObject:data waitUntilDone:YES];
    });
}

- (void)fetchedData:(NSData *)responseData
{
    if (responseData == nil)
    {
        NSLog(@"get User class result in nil");
        return;
    }
    
    //parse out the json data
    NSError* error;
    NSDictionary* json = [NSJSONSerialization JSONObjectWithData:responseData
                                                         options:kNilOptions
                                                           error:&error];
    NSLog(@"dictionary data %@",json);
    
    [[MeViewData singleton] reset];
    [[MeViewData singleton] setJsonDictionaryData:json];
    
    [classesTableView_ reloadData];
}


#pragma mark -
#pragma mark Table View Data Source Methods
-(NSInteger)tableView:(UITableView *)tableView numberOfRowsInSection:(NSInteger)section
{
    int count = [[MeViewData singleton] getNotificationCount];
    return count;
}

-(UITableViewCell *)tableView:(UITableView *)tableView cellForRowAtIndexPath:(NSIndexPath *)indexPath
{
    static NSString *cellIdentifier = @"cellIdentifier";
    
    UITableViewCell *cell = [tableView dequeueReusableCellWithIdentifier:cellIdentifier];
    if (cell == nil)
    {
        cell = [self getTableViewCellbyLoadingNib];
    }
    
    NSUInteger row = [indexPath row];
    
    NSMutableDictionary* dataItemEelements =  [[MeViewData singleton] getNotificationItem:(int)row];
    AdminPublishTableViewCell* adminPublishTableViewCell = (AdminPublishTableViewCell*)cell;
    
    NSString* className = [NotificationDataItem getClassName:dataItemEelements];
    NSString* classTime = [NotificationDataItem getClassTime:dataItemEelements];
    if (classTime == nil)
    {
        classTime = @"";
    }
    
    [adminPublishTableViewCell setClassName:className];
    [adminPublishTableViewCell setClassTime:classTime];
    return cell;
}

- (UITableViewCell*)getTableViewCellbyLoadingNib
{
    NSString* personTableViewCell = @"AdminPublishTableViewCell";
    
    AdminPublishTableViewCell *cell;
    NSArray *nib = [[NSBundle mainBundle] loadNibNamed:personTableViewCell
                                                 owner:self
                                               options:nil];
    for (id oneObject in nib)
    {
        if ([oneObject isKindOfClass:[AdminPublishTableViewCell class]])
        {
            cell = (AdminPublishTableViewCell *)oneObject;
        }
    }
    return cell;
}

-(CGFloat)tableView:(UITableView*)tableView heightForRowAtIndexPath:(NSIndexPath*)indexPath
{
    return [self getTableViewCellbyLoadingNib].frame.size.height;
}

#pragma mark -
#pragma mark Table Delegate Methods
-(void)tableView:(UITableView*)tableView didSelectRowAtIndexPath:(NSIndexPath*)indexPath
{
    NSUInteger row = [indexPath row];
    NSMutableDictionary* dataItem =  [[MeViewData singleton] getNotificationItem:(int)row];
    if (dataItem == nil)
    {
        return;
    }
    
    NotificationDataItem* notificationDataItem =  [[NotificationDataItem alloc] init];
    [notificationDataItem setDataItemEelements:dataItem];
    
    NotificationPublishViewController* notificationPublishViewController = [[NotificationPublishViewController alloc] init];
    [notificationPublishViewController setEditType:EditType_View];
    [notificationPublishViewController setNotificationDataItem:notificationDataItem];
    [self.navigationController pushViewController:notificationPublishViewController animated:YES];
}

@end
