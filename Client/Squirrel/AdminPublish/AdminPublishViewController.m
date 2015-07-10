//
//  AdminPublishViewController.m
//  Squirrel
//
//  Created by JamesMao on 4/11/15.
//  Copyright (c) 2015 JamesMao. All rights reserved.
//

#import "AdminPublishViewController.h"
#import "NotificationPublishViewController.h"
#import "AdminPublishTableViewCell.h"
#import "DataModel/UserViewData.h"
#import "Define.h"

@interface AdminPublishViewController ()
{
    IBOutlet UITableView* notificationTableView_;
}
@end

@implementation AdminPublishViewController

- (void)viewDidLoad
{
    [super viewDidLoad];

    self.navigationItem.rightBarButtonItem = [[UIBarButtonItem alloc] initWithTitle:@"+"
                                                                              style:UIBarButtonItemStyleDone
                                                                             target:self
                                                                             action:@selector(onAdd)];
    
    self.navigationController.navigationBar.titleTextAttributes = @{NSForegroundColorAttributeName: [UIColor whiteColor]};

}

- (void)viewWillAppear:(BOOL)animated
{
    [notificationTableView_ reloadData];
    
    NSString* classListURLString = [NSString stringWithFormat:@"%@%@", SERVER_IP, URL_PART_CLASS_LIST];
    [self request:RequestTypeGet urlString:classListURLString parameters:nil];
}

- (void)onSuccess:(AFHTTPRequestOperation *)operation responseObject:(id)responseObject
{
    [super onSuccess:operation responseObject:responseObject];
    
    [[UserViewData singleton] reset];
    [[UserViewData singleton] setJsonDictionaryData:responseObject];
    
    [notificationTableView_ reloadData];
}

- (void)onFail:(AFHTTPRequestOperation *)operation error:(NSError*)error
{
    [super onFail:operation error:error];
    
    [[UserViewData singleton] reset];    
    [notificationTableView_ reloadData];
}

- (void)onAdd
{
    NotificationPublishViewController* notificationPublishViewController = [[NotificationPublishViewController alloc] init];
    [notificationPublishViewController setEditType:EditType_Add];
    [self.navigationController pushViewController:notificationPublishViewController animated:YES];
}

- (void)didReceiveMemoryWarning
{
    [super didReceiveMemoryWarning];
    // Dispose of any resources that can be recreated.
}

#pragma mark -
#pragma mark Table View Data Source Methods
-(NSInteger)tableView:(UITableView *)tableView numberOfRowsInSection:(NSInteger)section
{
    int count = [[UserViewData singleton] getNotificationCount];
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
    
    NSMutableDictionary* dataItemEelements =  [[UserViewData singleton] getNotificationItem:(int)row];
    AdminPublishTableViewCell* adminPublishTableViewCell = (AdminPublishTableViewCell*)cell;
    
    NSString* className = [NotificationDataItem getClassName:dataItemEelements];
    NSString* classTime = [NotificationDataItem getClassTime:dataItemEelements];
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
    NSMutableDictionary* dataItem =  [[UserViewData singleton] getNotificationItem:(int)row];
    if (dataItem == nil)
    {
        return;
    }
    
    NotificationDataItem* notificationDataItem =  [[NotificationDataItem alloc] init];
    [notificationDataItem setDataItemEelements:dataItem];
    
    NotificationPublishViewController* notificationPublishViewController = [[NotificationPublishViewController alloc] init];
    [notificationPublishViewController setEditType:EditType_Editable];
    [notificationPublishViewController setNotificationDataItem:notificationDataItem];
    [self.navigationController pushViewController:notificationPublishViewController animated:YES];
}


@end
